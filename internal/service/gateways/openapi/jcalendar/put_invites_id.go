package jcalendar

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	cmdinvite "jcalendar/internal/service/usecase/commands/invite"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

// nolint:revive,stylecheck
func (s *Server) PutInvitesId(c echo.Context, id string) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t convert string id param for putting events: %v", err)
		return echo.ErrBadRequest
	}

	gcmd, err := qrinvite.NewGetInviteQuery(ctx, uint(iid))
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t create GetInvite query: %v", err)
		return echo.ErrBadRequest
	}

	inv, err := s.application.Queries.GetInvite.Handle(ctx, gcmd)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t execute GetInvite query: %v", err)
		if inv == nil {
			return echo.ErrInternalServerError
		}
		return echo.ErrNotFound
	}

	if !isResourceOwner(ctx, inv.UserID, c.Get(userIDClaim).(uint)) {
		return echo.ErrForbidden
	}

	req := &jcalendarsrv.InviteUpdateRequest{}
	if err = c.Bind(req); err != nil {
		logrus.WithContext(ctx).Errorf("can`t binds inviteUpdate request body: %v", err)
		return echo.ErrBadRequest
	}

	ucmd, err := cmdinvite.NewUpdateInviteCommand(ctx, uint(iid), *req.Data.IsAccepted)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t create UpdateInvite command: %v", err)
		return echo.ErrBadRequest
	}

	err = s.application.Commands.UpdateInvite.Handle(ctx, ucmd)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t execute UpdateInvite command: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(200, jcalendarsrv.UpdatedInvite{ID: &iid})
}

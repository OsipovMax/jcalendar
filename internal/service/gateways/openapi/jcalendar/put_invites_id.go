package jcalendar

import (
	"net/http"
	"strconv"

	cmdinvite "jcalendar/internal/service/usecase/commands/invite"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"

	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) PutInvitesId(c echo.Context, id string) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	gcmd, err := qrinvite.NewGetInviteQuery(ctx, uint(iid))
	if err != nil {
		return err
	}

	inv, err := s.application.Queries.GetInvite.Handle(ctx, gcmd)
	if err != nil {
		return err
	}

	if !isResourceOwner(ctx, inv.UserID, c.Get(userIDClaim).(uint)) {
		return echo.ErrForbidden
	}

	req := &jcalendarsrv.InviteUpdateRequest{}
	if err = c.Bind(req); err != nil {
		return err
	}

	ucmd, err := cmdinvite.NewUpdateInviteCommand(ctx, uint(iid), *req.Data.IsAccepted)
	if err != nil {
		return err
	}

	_, err = s.application.Commands.UpdateInvite.Handle(ctx, ucmd)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

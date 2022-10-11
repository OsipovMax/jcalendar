package jcalendar

import (
	"net/http"
	"strconv"

	"jcalendar/internal/service/usecase/commands/invite"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) PutInvitesId(c echo.Context, id string) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	req := &jcalendarsrv.InviteUpdateRequest{}
	if err = c.Bind(req); err != nil {
		return err
	}

	cmd, err := invite.NewUpdateInviteCommand(ctx, uint(iid), *req.IsAccepted) //TODO req.Data.IsAccepted
	if err != nil {
		return err
	}

	_, err = s.application.Commands.UpdateInvite.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

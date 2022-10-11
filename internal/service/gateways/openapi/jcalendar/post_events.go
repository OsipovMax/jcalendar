package jcalendar

import (
	"net/http"

	"jcalendar/internal/service/usecase/commands/event"
	"jcalendar/internal/service/usecase/commands/invite"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) PostEvents(c echo.Context) error {
	ctx := c.Request().Context()

	req := &jcalendarsrv.EventRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	cmd, err := event.NewCreateEventCommand(
		ctx,
		*req.Data.From,
		*req.Data.Till,
		uint(*req.Data.CreatorID),
		*req.Data.Participants,
		*req.Data.Details,
		*req.Data.IsPrivate,
		*req.Data.IsRepeat,
	)
	if err != nil {
		return err
	}

	eID, err := s.application.Commands.CreateEvent.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	for _, participantID := range cmd.ParticipantsIDs {
		var icmd *invite.CreateInviteCommand
		icmd, err = invite.NewCreateInviteCommand(ctx, participantID, eID, false)
		if err != nil {
			return err
		}

		_, err = s.application.Commands.CreateInvite.Handle(ctx, icmd)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusCreated, jcalendarsrv.CreatedEvent{ID: pcaster(int(eID))})
}

package jcalendar

import (
	"fmt"
	"net/http"

	"jcalendar/internal/service/usecase/commands/event"
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
		*req.Data.ScheduleRule,
		*req.Data.Details,
		*req.Data.IsPrivate,
		*req.Data.IsRepeat,
	)
	if err != nil {
		return err
	}

	eID, err := s.application.Commands.CreateEvent.Handle(ctx, cmd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.JSON(http.StatusCreated, jcalendarsrv.CreatedEvent{ID: pcaster(int(eID))})
}

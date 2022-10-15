package jcalendar

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"jcalendar/internal/service/usecase/commands/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func (s *Server) PostEvents(c echo.Context) error {
	ctx := c.Request().Context()

	req := &jcalendarsrv.EventRequest{}
	if err := c.Bind(req); err != nil {
		logrus.WithContext(ctx).Errorf("can`t binds event request body: %v", err)
		return echo.ErrBadRequest
	}

	cmd, err := event.NewCreateEventCommand(
		ctx,
		req.Data.From,
		req.Data.Till,
		uint(req.Data.CreatorID),
		req.Data.Participants,
		req.Data.ScheduleRule,
		req.Data.Details,
		req.Data.IsPrivate,
		req.Data.IsRepeat,
	)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t create new createEventCommand: %v", err)
		return echo.ErrBadRequest
	}

	eID, err := s.application.Commands.CreateEvent.Handle(ctx, cmd)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t execute createEventCommand: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, jcalendarsrv.CreatedEvent{ID: int(eID)})
}

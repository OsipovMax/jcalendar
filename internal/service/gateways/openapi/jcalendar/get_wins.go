package jcalendar

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	mevent "jcalendar/internal/service/usecase/managers/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func (s *Server) GetWindows(c echo.Context, params jcalendarsrv.GetWindowsParams) error {
	ctx := c.Request().Context()

	from, till, err := s.application.EventManager.GetClosestFreeWindow(ctx, params.UserIds, params.WinSize)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t get closest free window form event manager: %v", err)
		if errors.Is(err, mevent.ErrInvalidExecutingQuery) {
			return echo.ErrInternalServerError
		}
		return echo.ErrBadRequest
	}

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.FreeWindowResponse{
			Data: jcalendarsrv.FreeWindow{
				From: from.Format(time.RFC3339),
				Till: till.Format(time.RFC3339),
			},
		},
	)
}

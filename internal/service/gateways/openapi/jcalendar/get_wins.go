package jcalendar

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"jcalendar/internal/pkg"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func (s *Server) GetWindows(c echo.Context, params jcalendarsrv.GetWindowsParams) error {
	ctx := c.Request().Context()

	from, till, err := s.application.EventManager.GetClosestFreeWindow(ctx, params.UserIds, params.WinSize)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t get closest free window form event: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.FreeWindowResponse{
			Data: &jcalendarsrv.FreeWindow{
				From: pkg.Type2pointer(from.Format(time.RFC3339)),
				Till: pkg.Type2pointer(till.Format(time.RFC3339)),
			},
		},
	)
}

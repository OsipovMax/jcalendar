package jcalendar

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func (s *Server) GetWindows(c echo.Context, params jcalendarsrv.GetWindowsParams) error {
	ctx := c.Request().Context()

	from, till, err := s.application.Finder.GetClosestFreeWindow(ctx, params.UserIds, params.WinSize)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.FreeWindowResponse{
			Data: &jcalendarsrv.FreeWindow{
				From: pcaster(from.String()),
				Till: pcaster(till.String()),
			},
		},
	)
}

package jcalendar

import (
	"github.com/labstack/echo/v4"

	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func (s *Server) GetWindows(c echo.Context, params jcalendarsrv.GetWindowsParams) error {
	return nil
}

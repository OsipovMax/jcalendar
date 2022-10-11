package main

import (
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(jcalendar.GetConfirmedUserMiddleware())

	jcalendarsrv.RegisterHandlers(e,
		jcalendar.NewServer(nil),
	)
}

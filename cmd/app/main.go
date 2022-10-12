package main

import (
	"context"

	"jcalendar/internal/service/app"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := echo.New()
	e.Use(jcalendar.GetConfirmedUserMiddleware())

	application, err := app.NewApplication(ctx)
	if err != nil {
		return //TODO log.fatal
	}

	jcalendarsrv.RegisterHandlers(e,
		jcalendar.NewServer(application),
	)
}

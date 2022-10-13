package main

import (
	"context"
	"log"

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
		log.Fatalf("can`t create new application: %v", err)
	}

	jcalendarsrv.RegisterHandlers(e,
		jcalendar.NewServer(application),
	)

	if err = e.Start(":8080"); err != nil {
		log.Fatalf("can`t start http server: %v", err)
	}
}

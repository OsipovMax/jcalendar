package main

import (
	"context"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pkg.NewDB(ctx)
	if err != nil {
		logrus.WithContext(ctx).Fatalf("can`t get gorm db defenition: %v", err)
	}

	e := echo.New()
	e.Use(jcalendar.GetConfirmedUserMiddleware())

	application, err := app.NewApplication(ctx, db)
	if err != nil {
		logrus.WithContext(ctx).Fatalf("can`t create new application: %v", err)
	}

	jcalendarsrv.RegisterHandlers(e,
		jcalendar.NewServer(application),
	)

	if err = e.Start(":8080"); err != nil {
	}
}

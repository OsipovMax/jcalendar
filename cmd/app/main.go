package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pkg.NewDB(ctx)
	if err != nil {
		logrus.WithContext(ctx).Fatalf("can`t get gorm db defenition: %v", err)
	}

	application, err := app.NewApplication(ctx, db)
	if err != nil {
		logrus.WithContext(ctx).Fatalf("can`t create new application: %v", err)
	}

	e := echo.New()
	e.Use(jcalendar.GetConfirmedUserMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	jcalendarsrv.RegisterHandlersWithBaseURL(e, jcalendar.NewServer(application), "/api")

	if err = e.Start(":8080"); err != nil {
		logrus.WithContext(ctx).Fatalf("can`t start HTTP Server: %v", err)
	}
}

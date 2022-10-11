package jcalendar

import (
	"jcalendar/internal/service/app"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

type Server struct {
	application *app.Application
}

var _ jcalendarsrv.ServerInterface = (*Server)(nil)

func NewServer(application *app.Application) *Server {
	return &Server{
		application: application,
	}
}

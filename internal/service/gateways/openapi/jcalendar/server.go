package jcalendar

import (
	"context"

	euser "jcalendar/internal/service/entity/user"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

type Manager interface {
	GetUserByEmail(ctx context.Context, email string) (*euser.User, error)
	GetUserByID(ctx context.Context, id string) (*euser.User, error)
	CreateUser(ctx context.Context, user *euser.User) (uint, error)
}

type Server struct {
	manager Manager
}

var _ jcalendarsrv.ServerInterface = (*Server)(nil)

func NewServer(manager Manager) *Server {
	return &Server{
		manager: manager,
	}
}

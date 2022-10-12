package jcalendar

import (
	"net/http"

	"jcalendar/internal/service/usecase/commands/user"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) PostUsers(c echo.Context) error {
	ctx := c.Request().Context()

	req := &jcalendarsrv.UserRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	cmd, err := user.NewCreateUserCommand(
		ctx,
		*req.Data.FirstName,
		*req.Data.LastName,
		*req.Data.Email,
		calcHash(*req.Data.Password),
		*req.Data.TimeZoneOffset,
	)
	if err != nil {
		return err
	}

	uID, err := s.application.Commands.CreateUser.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, jcalendarsrv.CreatedUser{ID: pcaster(int(uID))})
}

package jcalendar

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	cmduser "jcalendar/internal/service/usecase/commands/user"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func (s *Server) PostUsers(c echo.Context) error {
	ctx := c.Request().Context()

	req := &jcalendarsrv.UserRequest{}
	if err := c.Bind(req); err != nil {
		logrus.WithContext(ctx).Errorf("can`t binds user request body: %v", err)
		return echo.ErrBadRequest
	}

	cmd, err := cmduser.NewCreateUserCommand(
		ctx,
		req.Data.FirstName,
		req.Data.LastName,
		req.Data.Email,
		calcHash(req.Data.Password),
		req.Data.TimeZoneOffset,
	)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t create CreateUser command: %v", err)
		return echo.ErrBadRequest
	}

	uID, err := s.application.Commands.CreateUser.Handle(ctx, cmd)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t execute CreateUser command: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, jcalendarsrv.CreatedUser{ID: int(uID)})
}

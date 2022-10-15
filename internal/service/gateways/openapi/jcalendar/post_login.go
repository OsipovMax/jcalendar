package jcalendar

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	qruser "jcalendar/internal/service/usecase/queries/user"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

const (
	usernameKey = "username"
	passwordKey = "password"
)

func (s *Server) PostLogin(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		username = c.FormValue(usernameKey)
		password = c.FormValue(passwordKey)
	)

	query, err := qruser.NewUserByEmailQuery(username)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t create UserByEmail query: %v", err)
		return echo.ErrBadRequest
	}

	u, err := s.application.Queries.GetUserByEmail.Handle(ctx, query)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t execute UserByEmail query: %v", err)
		if u == nil {
			return echo.ErrInternalServerError
		}

		return echo.ErrNotFound
	}

	if !isPasswordValid(u.HashedPassword, password) {
		return echo.ErrUnauthorized
	}

	st, err := generateJWT(ctx, u.ID, username)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t generate new JWT token: %v", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, jcalendarsrv.TokenResponse{Data: &jcalendarsrv.TokenPayload{Token: &st}})
}

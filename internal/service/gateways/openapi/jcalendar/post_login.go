package jcalendar

import (
	"net/http"
	"time"

	qruser "jcalendar/internal/service/usecase/queries/user"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	usernameKey = "username"
	passwordKey = "password"

	tokenLifetime = 10 * time.Minute
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	userID    uint
	userEmail string
	jwt.StandardClaims
}

func (s *Server) PostLogin(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		username = c.FormValue(usernameKey)
		password = c.FormValue(passwordKey)
	)

	query, err := qruser.NewUserByEmailQuery(username)
	if err != nil {
		return err
	}

	u, err := s.application.Queries.GetUserByEmail.Handle(ctx, query)
	if err != nil {
		return err
	}

	if u == nil || !isPasswordValid(u.HashedPassword, password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

	st, err := getJWT(ctx, u.ID, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Please retry later")
	}

	return c.JSON(http.StatusOK, jcalendarsrv.TokenResponse{Data: &jcalendarsrv.TokenPayload{Token: &st}})
}

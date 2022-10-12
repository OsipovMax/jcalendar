package jcalendar

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	AuthHeaderKey = "Authorization"
)

func MiddlewareSkipper(c echo.Context) bool {
	return (c.Request().URL.Path == "/users" || c.Request().URL.Path == "/login") && c.Request().Method == http.MethodPost
}

func GetConfirmedUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if MiddlewareSkipper(c) {
				return next(c)
			}

			parts := strings.Split(c.Request().Header.Get(AuthHeaderKey), " ")
			if len(parts) < 2 {
				return echo.NewHTTPError(http.StatusForbidden)
			}

			var (
				token = parts[1]
				ctx   = c.Request().Context()
			)

			userID, err := getUserID(ctx, token)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden)
			}

			c.Set(userIDClaim, userID)

			return next(c)
		}
	}
}

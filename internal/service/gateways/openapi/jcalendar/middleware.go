package jcalendar

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"jcalendar/internal/metrics"
)

const (
	AuthHeaderKey = "Authorization"
)

func MiddlewareSkipper(c echo.Context) bool {
	return c.Request().URL.Path == "/api/users" || c.Request().URL.Path == "/api/login" || c.Request().URL.Path == "/metrics"
}

func GetConfirmedUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			defer func() {
				var (
					statusCode = "2xx"
					eerr       *echo.HTTPError
				)

				if errors.As(err, &eerr) {
					statusCode = strconv.Itoa(eerr.Code)
				}

				metrics.HTTPRequestsTotal.WithLabelValues(c.Request().URL.String(), statusCode)
			}()

			if MiddlewareSkipper(c) {
				err = next(c)
				return err
			}

			parts := strings.Split(c.Request().Header.Get(AuthHeaderKey), " ")
			if len(parts) < 2 {
				return echo.ErrForbidden
			}

			var (
				token = parts[1]
				ctx   = c.Request().Context()
			)

			userID, err := getUserID(ctx, token)
			if err != nil {
				return echo.ErrForbidden
			}

			c.Set(userIDClaim, userID)

			err = next(c)
			return err
		}
	}
}

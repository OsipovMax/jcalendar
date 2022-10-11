package jcalendar

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	AuthHeaderKey = "Authorization"
)

func GetConfirmedUserMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*jwtCustomClaims)
			c.Set("userID", claims.userID)

			return next(c)
		}
	}
}

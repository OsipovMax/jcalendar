package jcalendar

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO:
func pcaster[T uint | string](val T) *T {
	return &val
}

func isPasswordValid(hashedPass, pass string) bool {
	return hashedPass == calcHash(pass)
}

func calcHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))

	return string(h.Sum(nil))
}

func getJWT(_ context.Context, uID uint, username string) (string, error) {
	claims := &jwtCustomClaims{
		uID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	st, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return st, nil
}

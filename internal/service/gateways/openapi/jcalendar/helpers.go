package jcalendar

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

const (
	userIDClaim   = "userID"
	userNameClaim = "userName"
	secretKey     = "secret" // Use vault for storage and receive at the start of the service !!!
)

func pcaster[T comparable](val T) *T {
	return &val
}

func isPasswordValid(hashedPass, pass string) bool {
	return hashedPass == calcHash(pass)
}

func calcHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))

	return hex.EncodeToString(h.Sum(nil))
}

func generateJWT(_ context.Context, uID uint, username string) (string, error) {
	//claims := &jwtCustomClaims{
	//	uID,
	//	username,
	//	jwt.StandardClaims{
	//		ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
	//	},
	//}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			userIDClaim:   uID,
			userNameClaim: username,
		},
	)

	st, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return st, nil
}

func getUserID(_ context.Context, token string) (uint, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		token, claims,
		func(token *jwt.Token) (interface{}, error) { return []byte(secretKey), nil },
	)
	if err != nil {
		return 0, fmt.Errorf("invalid jwtToken parsing: %w", err)
	}

	for k, v := range claims {
		if k == userIDClaim {
			return uint(v.(float64)), nil
		}
	}

	return 0, errors.New("missing userID claim")
}

func isResourceOwner(_ context.Context, resourceOwnerID, userID uint) bool {
	return resourceOwnerID == userID
}

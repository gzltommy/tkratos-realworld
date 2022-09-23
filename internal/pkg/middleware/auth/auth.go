package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var currentUserKey struct{}

type CurrentUser struct {
	UserID uint
}

func GenerateToken(secret string, userid uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
		"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func FromContext(ctx context.Context) *CurrentUser {
	return ctx.Value(currentUserKey).(*CurrentUser)
}

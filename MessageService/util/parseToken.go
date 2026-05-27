package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
)

func ParseToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	profileIDStr, ok := claims["profile_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("no profile_id in token")
	}

	return uuid.Parse(profileIDStr)
}

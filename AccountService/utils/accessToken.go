package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

func getJWTSecret() []byte { return []byte(os.Getenv("JWT_SECRET")) }

func getAccessTTL() time.Duration {
	minutesStr := os.Getenv("JWT_ACCESS_EXPIRES_MINUTES")
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		minutes = 15
	}
	return time.Duration(minutes) * time.Minute
}

func GenerateAccessToken(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"profile_id": id.String(),
		"exp":        time.Now().Add(getAccessTTL()).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

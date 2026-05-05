package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
)

func getRefreshTTL() time.Duration {
	daysStr := os.Getenv("JWT_REFRESH_EXPIRES_DAYS")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 15
	}
	return time.Duration(days) * 24 * time.Hour
}

func GenerateRefreshToken(profileID, refreshID uuid.UUID) (time.Time, string, error) {
	expTime := time.Now().Add(getRefreshTTL())
	claims := jwt.MapClaims{
		"profile_id": profileID.String(),
		"token_id":   refreshID.String(),
		"exp":        expTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(getJWTSecret())
	return expTime, tokenStr, err
}

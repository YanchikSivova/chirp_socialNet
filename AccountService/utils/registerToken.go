package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

func getSecret() []byte {
	return []byte(os.Getenv("JWT_REGISTER_SECRET"))
}

func getRegisterTTL() time.Duration {
	minutesStr := os.Getenv("JWT_REGISTER_EXPIRES_MINUTES")
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		minutes = 5
	}
	return time.Duration(minutes) * time.Minute
}

func GenerateRegisterToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(getRegisterTTL()).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

func ParseRegisterToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

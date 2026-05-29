package util

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestParseTokenReturnsProfileID(t *testing.T) {
	t.Setenv("JWT_SECRET", "message-secret")

	profileID := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"profile_id": profileID.String(),
		"exp":        time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("message-secret"))
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	got, err := ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}
	if got != profileID {
		t.Fatalf("profileID = %s, want %s", got, profileID)
	}
}

func TestParseTokenRejectsMissingProfileID(t *testing.T) {
	t.Setenv("JWT_SECRET", "message-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("message-secret"))
	if err != nil {
		t.Fatalf("SignedString returned error: %v", err)
	}

	if _, err := ParseToken(tokenStr); err == nil {
		t.Fatal("ParseToken accepted token without profile_id")
	}
}

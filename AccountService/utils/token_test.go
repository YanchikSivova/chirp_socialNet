package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAccessTokenIncludesProfileID(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_ACCESS_EXPIRES_MINUTES", "10")

	profileID := uuid.New()
	tokenStr, err := GenerateAccessToken(profileID)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("jwt.Parse returned error: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		t.Fatal("generated token is not valid")
	}
	if got := claims["profile_id"]; got != profileID.String() {
		t.Fatalf("profile_id = %v, want %s", got, profileID.String())
	}
	if _, ok := claims["iat"].(float64); !ok {
		t.Fatal("iat claim is missing")
	}
	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		t.Fatalf("exp claim is invalid: %v", claims["exp"])
	}
}

func TestGenerateTemporaryTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_TEMPORARY_SECRET", "temporary-secret")
	t.Setenv("JWT_TEMPORARY_EXPIRES_MINUTES", "5")

	const email = "user@example.com"
	token, err := GenerateTemporaryToken(email)
	if err != nil {
		t.Fatalf("GenerateTemporaryToken returned error: %v", err)
	}

	claims, err := ParseTemporaryToken(token)
	if err != nil {
		t.Fatalf("ParseTemporaryToken returned error: %v", err)
	}
	if got := claims["email"]; got != email {
		t.Fatalf("email = %v, want %s", got, email)
	}
}

func TestGenerateCodeReturnsSixDigits(t *testing.T) {
	code := GenerateCode()
	if len(code) != 6 {
		t.Fatalf("len(code) = %d, want 6", len(code))
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			t.Fatalf("code contains non-digit rune %q", r)
		}
	}
}

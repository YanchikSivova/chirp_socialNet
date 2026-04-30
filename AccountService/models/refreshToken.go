package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	RefreshTokenID uuid.UUID `json:"refresh_token_id"`
	ProfileID      uuid.UUID `json:"profile_id"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"expires_at"`
}

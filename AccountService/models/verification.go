package models

import (
	"github.com/google/uuid"
	"time"
)

type Verification struct {
	VerificationID uuid.UUID `json:"verification_id"`
	CredentialsID  uuid.UUID `json:"creentials_id"`
	Code           string    `json:"code"`
	ExpiresAT      time.Time `json:"expires_at"`
}

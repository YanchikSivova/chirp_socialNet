package models

import (
	"github.com/google/uuid"
	"time"
)

type Credentials struct {
	CredentialsID  uuid.UUID `json:"credentials_id"`
	ProfileID      uuid.UUID `json:"profile_id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	Status         string    `json:"status"`
}

package models

import "github.com/google/uuid"

type EmailChange struct {
	EmailChangeID uuid.UUID `json:"email_change_id"`
	CredentialsID uuid.UUID `json:"credentials_id"`
	NewEmail      string    `json:"new_email"`
}

package models

import "github.com/google/uuid"

type PasswordChange struct {
	PasswordChangeID uuid.UUID `json:"password_change_id"`
	CredentialsID    uuid.UUID `json:"credentials_id"`
	NewPasswordHash  string    `json:"new_password_hash"`
}

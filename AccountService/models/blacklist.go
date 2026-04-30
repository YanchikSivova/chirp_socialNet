package models

import "github.com/google/uuid"

type Blacklist struct {
	BlacklistID     uuid.UUID `json:"blacklist_id"`
	ProfileID       uuid.UUID `json:"profile_id"`
	BannedProfileID uuid.UUID `json:"banned_profile_id"`
}

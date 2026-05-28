package models

import "github.com/google/uuid"

type Profile struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}
type Profiles struct {
	Profiles []uuid.UUID `json:"profiles"`
}

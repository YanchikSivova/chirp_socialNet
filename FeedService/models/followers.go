package models

import "github.com/google/uuid"

type Followers struct {
	Followers []uuid.UUID `json:"followers"`
}

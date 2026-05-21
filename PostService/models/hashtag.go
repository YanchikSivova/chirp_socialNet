package models

import "github.com/google/uuid"

type Hashtag struct {
	HashtagID   uuid.UUID `json:"hashtag_id"`
	HashtagName string    `json:"hashtag_name"`
}

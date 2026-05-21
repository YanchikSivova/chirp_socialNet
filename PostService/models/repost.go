package models

import (
	"github.com/google/uuid"
	"time"
)

type Repost struct {
	RepostID  uuid.UUID `json:"repost_id"`
	PostID    uuid.UUID `json:"post_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	CreatedAt time.Time `json:"created_at"`
}

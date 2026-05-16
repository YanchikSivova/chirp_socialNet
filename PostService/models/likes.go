package models

import (
	"github.com/google/uuid"
	"time"
)

type Likes struct {
	LikesID   uuid.UUID `json:"likes_id"`
	PostID    uuid.UUID `json:"post_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	LikedAt   time.Time `json:"liked_at"`
}

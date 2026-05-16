package models

import (
	"github.com/google/uuid"
	"time"
)

type CommentLike struct {
	CommentLikeID uuid.UUID `json:"comment_like_id"`
	CommentID     uuid.UUID `json:"comment_id"`
	ProfileID     uuid.UUID `json:"profile_id"`
	LikedAt       time.Time `json:"liked_at"`
}

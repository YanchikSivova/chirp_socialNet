package kafkaEvents

import (
	"github.com/google/uuid"
	"time"
)

type PostEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	PostID    uuid.UUID `json:"post_id"`
	IsRepost  bool      `json:"is_repost"`
	CreatedAt time.Time `json:"created_at"`
}

package kafkaEvents

import "github.com/google/uuid"

type PostEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	PostID    uuid.UUID `json:"post_id"`
}

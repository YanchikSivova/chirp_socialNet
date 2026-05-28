package kafkaEvents

import (
	"github.com/google/uuid"
	"time"
)

type NotificationEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	ProfileID uuid.UUID `json:"profile_id"`
	ActorID   uuid.UUID `json:"actor_id"`
	Type      string    `json:"entity_type"`
	EntityID  string    `json:"entity_id"`
	CreatedAt time.Time `json:"created_at"`
}

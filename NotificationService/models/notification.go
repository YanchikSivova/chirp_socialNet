package models

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	NotificationID uuid.UUID `json:"notification_id"`
	ProfileID      uuid.UUID `json:"profile_id"`
	ActorID        uuid.UUID `json:"actor_id"`
	Type           string    `json:"entity_type"`
	EntityID       string    `json:"entity_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type NotificationsList struct {
	Notification Notification `json:"notification"`
	Actor        Profile      `json:"actor"`
}

package models

import "github.com/google/uuid"

type Subscription struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
	SubscriberID   uuid.UUID `json:"subscriber_id"`
	SubscribedID   uuid.UUID `json:"subscribed_id"`
}

type Followers struct {
	Followers []uuid.UUID `json:"followers"`
}

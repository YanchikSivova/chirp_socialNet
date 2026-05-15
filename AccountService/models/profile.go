package models

import "github.com/google/uuid"

type Profile struct {
	ProfileID         uuid.UUID `json:"profile_id"`
	Name              *string   `json:"name"`
	Username          *string   `json:"username"`
	Avatar            *string   `json:"avatar"`
	Description       *string   `json:"description"`
	SubscribersAmount int       `json:"subscribers_amount"`
	SubscribedAmount  int       `json:"subscribed_amount"`
	PostsAmount       int       `json:"posts_amount"`
	IsCompleted       bool      `json:"is_completed"`
}

type ProfileMinimum struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

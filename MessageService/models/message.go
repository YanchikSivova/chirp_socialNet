package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	MessageID uuid.UUID `json:"message_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsMine    bool      `json:"is_mine"`
}

type MessageResponse struct {
	MessageID uuid.UUID `json:"message_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsMine    bool      `json:"is_mine"`

	Post *PostPreview `json:"post"`
}

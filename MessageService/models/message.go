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

type WSMessage struct {
	Type           string    `json:"type"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Content        string    `json:"content"`
	ContentType    string    `json:"content_type"`
}

type WSMessageResponse struct {
	Type           string    `json:"type"`
	MessageID      uuid.UUID `json:"message_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Content        string    `json:"content"`
	ContentType    string    `json:"content_type"`
	CreatedAt      time.Time `json:"created_at"`
	ProfileID      uuid.UUID `json:"profile_id"`
}

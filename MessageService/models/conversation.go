package models

import (
	"github.com/google/uuid"
	"time"
)

type ConversationResponse struct {
	ConversationID uuid.UUID         `json:"conversation_id"`
	Profile        Profile           `json:"profile"`
	Messages       []MessageResponse `json:"messages"`
	NextCursor     *time.Time        `json:"next_cursor"`
}

type ConversationPreview struct {
	ConversationID    uuid.UUID  `json:"conversation_id"`
	Profile           Profile    `json:"profile"`
	LastReadMessageID *uuid.UUID `json:"last_read_message_id"`
	UnreadCount       int        `json:"unread_count"`
}

type ConversationPreviewList struct {
	Conversations []ConversationPreview `json:"conversations"`
}

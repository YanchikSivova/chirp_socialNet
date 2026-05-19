package models

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	CommentID       uuid.UUID  `json:"comment_id"`
	PostID          uuid.UUID  `json:"post_id"`
	ProfileID       uuid.UUID  `json:"profile_id"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id"`
	Content         string     `json:"content"`
	LikesAmount     int        `json:"likes_amount"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CommentResponse struct {
	Author      Profile   `json:"author"`
	CommentID   uuid.UUID `json:"comment_id"`
	Content     string    `json:"content"`
	LikesAmount int       `json:"likes_amount"`
	CreatedAt   time.Time `json:"created_at"`
	IsLiked     bool      `json:"is_liked"`
}

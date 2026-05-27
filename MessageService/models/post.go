package models

import "github.com/google/uuid"

type PostPreview struct {
	PostID  uuid.UUID `json:"post_id"`
	Content string    `json:"content"`
}

type PostPreviews struct {
	Posts []PostPreview `json:"posts"`
}

type PostIDs struct {
	Posts []uuid.UUID `json:"posts"`
}

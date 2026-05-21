package models

import "github.com/google/uuid"

type Image struct {
	ImageID    uuid.UUID `json:"image_id"`
	PostID     uuid.UUID `json:"post_id"`
	Url        string    `json:"url"`
	OrderIndex int       `json:"order_index"`
}

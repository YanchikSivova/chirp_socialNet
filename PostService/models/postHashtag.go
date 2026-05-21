package models

import "github.com/google/uuid"

type PostHashtag struct {
	PostHashtagID uuid.UUID `json:"post_hashtag_id"`
	PostID        uuid.UUID `json:"post_id"`
	HashtagID     uuid.UUID `json:"hashtag_id"`
	OrderIndex    int       `json:"order_index"`
}

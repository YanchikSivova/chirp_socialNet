package models

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	PostID         uuid.UUID  `json:"post_id"`
	ProfileID      uuid.UUID  `json:"profile_id"`
	Content        string     `json:"content"`
	PublishedAt    *time.Time `json:"published_at"`
	LastEditedAt   time.Time  `json:"last_edited_at"`
	Status         string     `json:"status"`
	LikesAmount    int        `json:"likes_amount"`
	CommentsAmount int        `json:"comments_amount"`
	RepostsAmount  int        `json:"reposts_amount"`
	ReportsAmount  int        `json:"reports_amount"`
}

type ImageList struct {
	Url        string `json:"url"`
	OrderIndex int    `json:"order_index"`
}

type HashtagList struct {
	HashtagName string `json:"hashtag_name"`
	OrderIndex  int    `json:"order_index"`
}
type CreatePostRequest struct {
	Content  string        `json:"content"`
	Images   []ImageList   `json:"images"`
	Hashtags []HashtagList `json:"hashtags"`
	Status   string        `json:"status"`
}

type UpdatePostRequest struct {
	Content  string        `json:"content"`
	Images   []ImageList   `json:"images"`
	Hashtags []HashtagList `json:"hashtags"`
}

type Profile struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

type PostResponse struct {
	Author         Profile       `json:"author"`
	PostID         uuid.UUID     `json:"post_id"`
	Content        string        `json:"content"`
	Images         []ImageList   `json:"images"`
	Hashtags       []HashtagList `json:"hashtags"`
	LikesAmount    int           `json:"likes_amount"`
	CommentsAmount int           `json:"comments_amount"`
	RepostsAmount  int           `json:"reposts_amount"`
	LastEditedAt   time.Time     `json:"last_edited_at"`
	PublishedAt    time.Time     `json:"published_at"`
	IsLiked        bool          `json:"is_liked"`
	IsReposted     bool          `json:"is_reposted"`
}

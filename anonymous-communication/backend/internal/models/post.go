package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ImageURL  string    `json:"image_url"`
	Caption   *string   `json:"caption,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PostAuthor struct {
	ID                string  `json:"id"`
	Username          string  `json:"username"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
}

type PostFeedItem struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	ImageURL  string     `json:"image_url"`
	Caption   *string    `json:"caption,omitempty"`
	LikeCount int64      `json:"like_count"`
	CreatedAt time.Time  `json:"created_at"`
	User      PostAuthor `json:"user"`
}

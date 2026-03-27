package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ImageURL  *string   `json:"image_url,omitempty"`
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
	ImageURL  *string    `json:"image_url,omitempty"`
	Caption   *string    `json:"caption,omitempty"`
	LikeCount int64      `json:"like_count"`
	CreatedAt time.Time  `json:"created_at"`
	User      PostAuthor `json:"user"`
}

type CreatePostUploadRequest struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	FileSize    int64  `json:"file_size"`
}

type CreatePostFromUploadRequest struct {
	ObjectPath string `json:"object_path"`
	Caption    string `json:"caption"`
}

type PostUploadTarget struct {
	ObjectPath string `json:"object_path"`
	UploadURL  string `json:"upload_url"`
	PublicURL  string `json:"public_url"`
	Method     string `json:"method"`
}

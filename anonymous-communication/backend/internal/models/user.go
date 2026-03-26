package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                        uuid.UUID `json:"id"`
	Username                  string    `json:"username"`
	Email                     string    `json:"email"`
	PasswordHash              string    `json:"-"`
	Role                      string    `json:"role"`
	ImpersonationPasswordHash *string   `json:"-"`
	ProfilePictureURL         *string   `json:"profile_picture_url,omitempty"`
	CreatedAt                 time.Time `json:"created_at"`
}

type UserResponse struct {
	ID                string    `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthSession struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:                u.ID.String(),
		Username:          u.Username,
		Email:             u.Email,
		Role:              u.Role,
		ProfilePictureURL: u.ProfilePictureURL,
		CreatedAt:         u.CreatedAt,
	}
}

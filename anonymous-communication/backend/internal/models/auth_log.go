package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthLog struct {
	ID            int64      `json:"id"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	Status        string     `json:"status"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     *string    `json:"user_agent,omitempty"`
	FailureReason *string    `json:"failure_reason,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateAuthLogParams struct {
	UserID        *uuid.UUID
	Status        string
	IPAddress     string
	UserAgent     string
	FailureReason string
}

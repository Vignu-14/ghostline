package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminAuditLog struct {
	ID           int64          `json:"id"`
	AdminID      uuid.UUID      `json:"admin_id"`
	TargetUserID *uuid.UUID     `json:"target_user_id,omitempty"`
	Action       string         `json:"action"`
	IPAddress    string         `json:"ip_address"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	DeleteModeMe              = "me"
	DeleteModeEveryone        = "everyone"
	DeletedMessagePlaceholder = "This message was deleted"
)

type Message struct {
	ID                   uuid.UUID  `json:"id"`
	SenderID             uuid.UUID  `json:"sender_id"`
	ReceiverID           uuid.UUID  `json:"receiver_id"`
	Content              string     `json:"content"`
	IsRead               bool       `json:"is_read"`
	CreatedAt            time.Time  `json:"created_at"`
	DeletedForSenderAt   *time.Time `json:"-"`
	DeletedForReceiverAt *time.Time `json:"-"`
	DeletedForEveryoneAt *time.Time `json:"-"`
}

type MessageResponse struct {
	ID                 string    `json:"id"`
	SenderID           string    `json:"sender_id"`
	ReceiverID         string    `json:"receiver_id"`
	Content            string    `json:"content"`
	IsRead             bool      `json:"is_read"`
	CreatedAt          time.Time `json:"created_at"`
	DeletedForEveryone bool      `json:"deleted_for_everyone"`
}

type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

type DeleteMessagesRequest struct {
	MessageIDs []string `json:"message_ids"`
	Mode       string   `json:"mode"`
}

type ClearConversationRequest struct {
	Mode string `json:"mode"`
}

type ConversationSummary struct {
	User        PostAuthor      `json:"user"`
	LastMessage MessageResponse `json:"last_message"`
	UnreadCount int64           `json:"unread_count"`
}

func (m Message) ToResponse() MessageResponse {
	content := m.Content
	deletedForEveryone := m.DeletedForEveryoneAt != nil
	if deletedForEveryone {
		content = DeletedMessagePlaceholder
	}

	return MessageResponse{
		ID:                 m.ID.String(),
		SenderID:           m.SenderID.String(),
		ReceiverID:         m.ReceiverID.String(),
		Content:            content,
		IsRead:             m.IsRead,
		CreatedAt:          m.CreatedAt,
		DeletedForEveryone: deletedForEveryone,
	}
}

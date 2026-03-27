package websocket

import "anonymous-communication/backend/internal/models"

const (
	EventTypeConnected = "connected"
	EventTypeMessage   = "message"
	EventTypeError     = "error"
)

type IncomingMessage struct {
	Type       string `json:"type"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

type Event struct {
	Type    string                  `json:"type"`
	Message *models.MessageResponse `json:"message,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

func NewConnectedEvent() Event {
	return Event{Type: EventTypeConnected}
}

func NewMessageEvent(message *models.MessageResponse) Event {
	return Event{
		Type:    EventTypeMessage,
		Message: message,
	}
}

func NewErrorEvent(message string) Event {
	return Event{
		Type:  EventTypeError,
		Error: message,
	}
}

package services

import (
	"context"
	"fmt"

	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type chatMessageRepository interface {
	Create(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*models.Message, error)
	Conversation(ctx context.Context, userID, otherUserID uuid.UUID, limit, offset int) ([]models.MessageResponse, error)
	ListConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.ConversationSummary, error)
	MarkConversationAsRead(ctx context.Context, userID, otherUserID uuid.UUID) error
}

type chatUserRepository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type ChatService struct {
	messages chatMessageRepository
	users    chatUserRepository
}

func NewChatService(messages chatMessageRepository, users chatUserRepository) *ChatService {
	return &ChatService{
		messages: messages,
		users:    users,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, senderID uuid.UUID, request models.SendMessageRequest) (*models.MessageResponse, error) {
	receiverID, err := utils.ParseUUID(request.ReceiverID)
	if err != nil {
		return nil, models.NewValidationError(map[string]string{
			"receiver_id": "receiver_id must be a valid uuid",
		})
	}

	if senderID == receiverID {
		return nil, models.ErrCannotMessageSelf
	}

	if _, err := s.users.FindByID(ctx, receiverID); err != nil {
		return nil, err
	}

	content := utils.SanitizeText(request.Content)
	if content == "" {
		return nil, models.NewValidationError(map[string]string{
			"content": "content is required",
		})
	}

	if len(content) > 5000 {
		return nil, models.NewValidationError(map[string]string{
			"content": "content must be 5000 characters or fewer",
		})
	}

	message, err := s.messages.Create(ctx, senderID, receiverID, content)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	response := message.ToResponse()
	return &response, nil
}

func (s *ChatService) GetConversation(ctx context.Context, userID, otherUserID uuid.UUID, limit, offset int) ([]models.MessageResponse, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	if _, err := s.users.FindByID(ctx, otherUserID); err != nil {
		return nil, err
	}

	if err := s.messages.MarkConversationAsRead(ctx, userID, otherUserID); err != nil {
		return nil, fmt.Errorf("mark messages as read: %w", err)
	}

	messages, err := s.messages.Conversation(ctx, userID, otherUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get conversation: %w", err)
	}

	return messages, nil
}

func (s *ChatService) ListConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.ConversationSummary, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	conversations, err := s.messages.ListConversations(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list conversations: %w", err)
	}

	return conversations, nil
}

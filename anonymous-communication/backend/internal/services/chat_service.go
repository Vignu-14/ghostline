package services

import (
	"context"
	"fmt"
	"strings"

	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/google/uuid"
)

type chatMessageRepository interface {
	Create(ctx context.Context, senderID, receiverID uuid.UUID, content string) (*models.Message, error)
	Conversation(ctx context.Context, userID, otherUserID uuid.UUID, limit, offset int) ([]models.MessageResponse, error)
	ListConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.ConversationSummary, error)
	MarkConversationAsRead(ctx context.Context, userID, otherUserID uuid.UUID) error
	FindByIDsForUser(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) ([]models.Message, error)
	DeleteForUser(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) (int64, error)
	DeleteForEveryone(ctx context.Context, userID uuid.UUID, messageIDs []uuid.UUID) (int64, error)
	ClearConversationForUser(ctx context.Context, userID, otherUserID uuid.UUID) (int64, error)
	DeleteConversationForEveryone(ctx context.Context, userID, otherUserID uuid.UUID) (int64, error)
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

func (s *ChatService) DeleteMessages(ctx context.Context, userID uuid.UUID, request models.DeleteMessagesRequest) (int64, error) {
	messageIDs, err := parseUniqueMessageIDs(request.MessageIDs)
	if err != nil {
		return 0, err
	}

	mode := strings.TrimSpace(strings.ToLower(request.Mode))
	if mode != models.DeleteModeMe && mode != models.DeleteModeEveryone {
		return 0, models.NewValidationError(map[string]string{
			"mode": "mode must be 'me' or 'everyone'",
		})
	}

	messages, err := s.messages.FindByIDsForUser(ctx, userID, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("load selected messages: %w", err)
	}
	if len(messages) != len(messageIDs) {
		return 0, models.ErrMessageNotFound
	}

	if mode == models.DeleteModeEveryone {
		sentMessageIDs := make([]uuid.UUID, 0, len(messages))
		receivedMessageIDs := make([]uuid.UUID, 0, len(messages))
		for _, message := range messages {
			if message.SenderID == userID {
				sentMessageIDs = append(sentMessageIDs, message.ID)
				continue
			}

			receivedMessageIDs = append(receivedMessageIDs, message.ID)
		}

		var affected int64
		if len(sentMessageIDs) > 0 {
			sentAffected, err := s.messages.DeleteForEveryone(ctx, userID, sentMessageIDs)
			if err != nil {
				return 0, fmt.Errorf("delete sent messages for everyone: %w", err)
			}
			affected += sentAffected
		}

		if len(receivedMessageIDs) > 0 {
			receivedAffected, err := s.messages.DeleteForUser(ctx, userID, receivedMessageIDs)
			if err != nil {
				return 0, fmt.Errorf("delete received messages for user: %w", err)
			}
			affected += receivedAffected
		}

		return affected, nil
	}

	affected, err := s.messages.DeleteForUser(ctx, userID, messageIDs)
	if err != nil {
		return 0, fmt.Errorf("delete messages for user: %w", err)
	}

	return affected, nil
}

func (s *ChatService) ClearConversation(ctx context.Context, userID, otherUserID uuid.UUID, mode string) (int64, error) {
	if userID == otherUserID {
		return 0, models.ErrCannotMessageSelf
	}

	if _, err := s.users.FindByID(ctx, otherUserID); err != nil {
		return 0, err
	}

	normalizedMode := strings.TrimSpace(strings.ToLower(mode))
	if normalizedMode == "" {
		normalizedMode = models.DeleteModeMe
	}
	if normalizedMode != models.DeleteModeMe && normalizedMode != models.DeleteModeEveryone {
		return 0, models.NewValidationError(map[string]string{
			"mode": "mode must be 'me' or 'everyone'",
		})
	}

	affected, err := s.messages.ClearConversationForUser(ctx, userID, otherUserID)
	if err != nil {
		return 0, fmt.Errorf("clear conversation: %w", err)
	}

	if normalizedMode == models.DeleteModeEveryone {
		sentAffected, err := s.messages.DeleteConversationForEveryone(ctx, userID, otherUserID)
		if err != nil {
			return 0, fmt.Errorf("delete sent conversation messages for everyone: %w", err)
		}

		affected += sentAffected
	}

	return affected, nil
}

func parseUniqueMessageIDs(input []string) ([]uuid.UUID, error) {
	if len(input) == 0 {
		return nil, models.NewValidationError(map[string]string{
			"message_ids": "at least one message id is required",
		})
	}

	seen := make(map[uuid.UUID]struct{}, len(input))
	ids := make([]uuid.UUID, 0, len(input))
	for _, rawID := range input {
		messageID, err := utils.ParseUUID(rawID)
		if err != nil {
			return nil, models.NewValidationError(map[string]string{
				"message_ids": "all message ids must be valid uuids",
			})
		}

		if _, exists := seen[messageID]; exists {
			continue
		}

		seen[messageID] = struct{}{}
		ids = append(ids, messageID)
	}

	return ids, nil
}

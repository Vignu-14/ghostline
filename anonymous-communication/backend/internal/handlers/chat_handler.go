package handlers

import (
	"errors"

	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) ListConversations(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	conversations, err := h.chatService.ListConversations(c.UserContext(), userID, limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return utils.Success(c, fiber.StatusOK, "conversations fetched successfully", fiber.Map{
		"conversations": conversations,
		"page":          page,
		"limit":         limit,
	})
}

func (h *ChatHandler) GetConversation(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	otherUserID, err := utils.ParseUUID(c.Params("userId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid user id", nil)
	}

	limit := c.QueryInt("limit", 50)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	messages, err := h.chatService.GetConversation(c.UserContext(), userID, otherUserID, limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "conversation fetched successfully", fiber.Map{
		"messages": messages,
		"page":     page,
		"limit":    limit,
	})
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	var request models.SendMessageRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	message, err := h.chatService.SendMessage(c.UserContext(), userID, request)
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusCreated, "message sent successfully", fiber.Map{
		"message": message,
	})
}

func (h *ChatHandler) DeleteMessages(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	var request models.DeleteMessagesRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	deletedCount, err := h.chatService.DeleteMessages(c.UserContext(), userID, request)
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "messages deleted successfully", fiber.Map{
		"deleted_count": deletedCount,
		"mode":          request.Mode,
	})
}

func (h *ChatHandler) ClearConversation(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	otherUserID, err := utils.ParseUUID(c.Params("userId"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid user id", nil)
	}

	var request models.ClearConversationRequest
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&request); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
		}
	}

	deletedCount, err := h.chatService.ClearConversation(c.UserContext(), userID, otherUserID, request.Mode)
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "conversation cleared successfully", fiber.Map{
		"deleted_count": deletedCount,
		"mode":          request.Mode,
	})
}

func (h *ChatHandler) handleError(c *fiber.Ctx, err error) error {
	var validationErr *models.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return utils.Error(c, fiber.StatusBadRequest, validationErr.Error(), validationErr.Fields)
	case errors.Is(err, models.ErrUserNotFound):
		return utils.Error(c, fiber.StatusNotFound, "user not found", nil)
	case errors.Is(err, models.ErrMessageNotFound):
		return utils.Error(c, fiber.StatusNotFound, "message not found", nil)
	case errors.Is(err, models.ErrCannotMessageSelf):
		return utils.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, models.ErrDeleteForEveryoneNotAllowed):
		return utils.Error(c, fiber.StatusForbidden, err.Error(), nil)
	default:
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}
}

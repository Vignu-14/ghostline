package handlers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"
	internalws "anonymous-communication/backend/internal/websocket"

	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketHandler struct {
	chatService *services.ChatService
	hub         *internalws.Hub
	jwtConfig   config.JWTConfig
	rateLimiter *middleware.RateLimiter
}

func NewWebSocketHandler(chatService *services.ChatService, hub *internalws.Hub, jwtConfig config.JWTConfig, rateLimiter *middleware.RateLimiter) *WebSocketHandler {
	return &WebSocketHandler{
		chatService: chatService,
		hub:         hub,
		jwtConfig:   jwtConfig,
		rateLimiter: rateLimiter,
	}
}

func (h *WebSocketHandler) Upgrade(c *fiber.Ctx) error {
	if !fiberws.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Try to get token from cookie first
	token := strings.TrimSpace(c.Cookies(h.jwtConfig.CookieName))
	
	// If no token in cookie, try query parameter
	if token == "" {
		token = strings.TrimSpace(c.Query("token"))
	}
	
	if token == "" {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	claims, err := utils.ParseToken(h.jwtConfig.Secret, token)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "invalid or expired session", nil)
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}

func (h *WebSocketHandler) HandleConnection(conn *fiberws.Conn) {
	userIDValue := conn.Locals("user_id")
	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		_ = conn.WriteJSON(internalws.NewErrorEvent("authentication required"))
		_ = conn.Close()
		return
	}

	client := internalws.NewClient(userID, conn)
	h.hub.Register(client)
	defer func() {
		h.hub.Unregister(client)
		_ = client.Close()
	}()

	_ = client.WriteJSON(internalws.NewConnectedEvent())

	for {
		var inbound internalws.IncomingMessage
		if err := conn.ReadJSON(&inbound); err != nil {
			break
		}

		messageType := inbound.Type
		if messageType == "" {
			messageType = internalws.EventTypeMessage
		}

		if messageType != internalws.EventTypeMessage {
			_ = client.WriteJSON(internalws.NewErrorEvent("unsupported websocket event type"))
			continue
		}

		senderID, err := utils.ParseUUID(userID)
		if err != nil {
			_ = client.WriteJSON(internalws.NewErrorEvent("invalid authenticated user"))
			continue
		}

		if h.rateLimiter != nil {
			allowed, retryAfter := h.rateLimiter.AllowMessageForUser(userID)
			if !allowed {
				retryAfterSeconds := int(math.Ceil(retryAfter.Seconds()))
				if retryAfterSeconds < 1 {
					retryAfterSeconds = 1
				}

				_ = client.WriteJSON(internalws.NewErrorEvent(fmt.Sprintf("message rate limit exceeded. try again in %d seconds.", retryAfterSeconds)))
				continue
			}
		}

		message, err := h.chatService.SendMessage(context.Background(), senderID, models.SendMessageRequest{
			ReceiverID: inbound.ReceiverID,
			Content:    inbound.Content,
		})
		if err != nil {
			_ = client.WriteJSON(internalws.NewErrorEvent(h.websocketErrorMessage(err)))
			continue
		}

		event := internalws.NewMessageEvent(message)
		_ = client.WriteJSON(event)
		internalws.BroadcastToUser(h.hub, message.ReceiverID, event)
	}
}

func (h *WebSocketHandler) websocketErrorMessage(err error) string {
	var validationErr *models.ValidationError

	switch {
	case errors.As(err, &validationErr):
		for _, message := range validationErr.Fields {
			return message
		}
		return validationErr.Error()
	case errors.Is(err, models.ErrUserNotFound):
		return "receiver not found"
	case errors.Is(err, models.ErrCannotMessageSelf):
		return err.Error()
	default:
		return "unable to process websocket message"
	}
}

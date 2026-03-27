package middleware

import (
	"strings"

	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	RequestIDKey    = "request_id"
	RequestIDHeader = "X-Request-ID"
)

func RequestID(c *fiber.Ctx) error {
	requestID := strings.TrimSpace(c.Get(RequestIDHeader))
	if requestID == "" {
		requestID = utils.NewUUID().String()
	}

	c.Locals(RequestIDKey, requestID)
	c.Set(RequestIDHeader, requestID)

	return c.Next()
}

func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Locals(RequestIDKey).(string)
	if ok {
		return requestID
	}

	return strings.TrimSpace(c.Get(RequestIDHeader))
}

package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(c *fiber.Ctx) error {
	startedAt := time.Now()
	err := c.Next()

	attributes := []any{
		"request_id", GetRequestID(c),
		"method", c.Method(),
		"path", c.OriginalURL(),
		"status", c.Response().StatusCode(),
		"latency_ms", time.Since(startedAt).Milliseconds(),
		"ip", c.IP(),
	}

	if claims, ok := GetClaims(c); ok {
		attributes = append(attributes, "user_id", claims.UserID)
		if claims.ImpersonatorID != nil {
			attributes = append(attributes, "impersonator_id", *claims.ImpersonatorID)
		}
	}

	switch {
	case err != nil:
		attributes = append(attributes, "error", err.Error())
		slog.Error("request failed", attributes...)
	case c.Response().StatusCode() >= fiber.StatusInternalServerError:
		slog.Error("request completed", attributes...)
	case c.Response().StatusCode() >= fiber.StatusBadRequest:
		slog.Warn("request completed", attributes...)
	default:
		slog.Info("request completed", attributes...)
	}

	return err
}

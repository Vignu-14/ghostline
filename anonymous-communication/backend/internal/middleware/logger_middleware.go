package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(c *fiber.Ctx) error {
	startedAt := time.Now()
	err := c.Next()
	duration := time.Since(startedAt)

	attributes := []any{
		"request_id", GetRequestID(c),
		"method", c.Method(),
		"path", c.OriginalURL(),
		"status", c.Response().StatusCode(),
		"latency_ms", duration.Milliseconds(),
		"ip", c.IP(),
		"origin", c.Get("Origin"),
	}

	// Add user info if authenticated
	if claims, ok := GetClaims(c); ok {
		attributes = append(attributes, "user_id", claims.UserID)
		if claims.ImpersonatorID != nil {
			attributes = append(attributes, "impersonator_id", *claims.ImpersonatorID)
		}
	}

	// Log auth endpoint failures with extra detail for debugging
	isAuthEndpoint := c.Path() == "/api/auth/login" || c.Path() == "/api/auth/register"
	statusCode := c.Response().StatusCode()

	switch {
	case err != nil:
		attributes = append(attributes, "error", err.Error())
		if isAuthEndpoint {
			slog.Error("auth request failed", attributes...)
		} else {
			slog.Error("request failed", attributes...)
		}
	case statusCode >= fiber.StatusInternalServerError:
		slog.Error("request completed", attributes...)
	case statusCode >= fiber.StatusBadRequest:
		if isAuthEndpoint {
			slog.Warn("auth request completed with error", attributes...)
		} else {
			slog.Warn("request completed", attributes...)
		}
	default:
		slog.Info("request completed", attributes...)
	}

	return err
}

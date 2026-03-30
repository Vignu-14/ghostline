package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewSecureHeaders(environment string) fiber.Handler {
	isProduction := strings.EqualFold(strings.TrimSpace(environment), "production")

	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Set("Cross-Origin-Opener-Policy", "same-origin")
		// Removed Cross-Origin-Resource-Policy to allow cross-origin API access
		// This header is not needed for API endpoints and can cause CORS issues
		c.Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")

		if isProduction {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		return c.Next()
	}
}

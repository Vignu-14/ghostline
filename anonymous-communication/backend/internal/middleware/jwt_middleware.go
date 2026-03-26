package middleware

import (
	"strings"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const AuthClaimsKey = "auth_claims"

type JWTMiddleware struct {
	config config.JWTConfig
}

func NewJWTMiddleware(cfg config.JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{config: cfg}
}

func (m *JWTMiddleware) RequireAuth(c *fiber.Ctx) error {
	token := strings.TrimSpace(c.Cookies(m.config.CookieName))
	if token == "" {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	claims, err := utils.ParseToken(m.config.Secret, token)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "invalid or expired session", nil)
	}

	c.Locals(AuthClaimsKey, claims)
	return c.Next()
}

func GetClaims(c *fiber.Ctx) (*models.JWTClaims, bool) {
	claims, ok := c.Locals(AuthClaimsKey).(*models.JWTClaims)
	return claims, ok
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, ok := GetClaims(c)
	if !ok {
		return uuid.Nil, models.ErrUnauthorized
	}

	return claims.UserUUID()
}

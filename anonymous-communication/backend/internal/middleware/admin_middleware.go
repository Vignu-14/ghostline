package middleware

import (
	"strings"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) RequireAdmin(c *fiber.Ctx) error {
	claims, ok := GetClaims(c)
	if !ok {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	if claims.ImpersonatorID != nil {
		return utils.Error(c, fiber.StatusForbidden, "impersonated sessions cannot access admin routes", nil)
	}

	if !strings.EqualFold(strings.TrimSpace(claims.Role), config.RoleAdmin) {
		return utils.Error(c, fiber.StatusForbidden, "admin access required", nil)
	}

	return c.Next()
}

func GetImpersonatorID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, ok := GetClaims(c)
	if !ok {
		return uuid.Nil, models.ErrUnauthorized
	}

	if claims.ImpersonatorID == nil {
		return uuid.Nil, models.ErrImpersonationNotActive
	}

	impersonatorID, err := uuid.Parse(*claims.ImpersonatorID)
	if err != nil {
		return uuid.Nil, err
	}

	return impersonatorID, nil
}

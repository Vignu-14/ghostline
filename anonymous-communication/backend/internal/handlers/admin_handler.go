package handlers

import (
	"errors"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	impersonationService *services.ImpersonationService
	jwtConfig            config.JWTConfig
}

func NewAdminHandler(impersonationService *services.ImpersonationService, jwtConfig config.JWTConfig) *AdminHandler {
	return &AdminHandler{
		impersonationService: impersonationService,
		jwtConfig:            jwtConfig,
	}
}

func (h *AdminHandler) Impersonate(c *fiber.Ctx) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	var request models.ImpersonateRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	session, err := h.impersonationService.Start(c.UserContext(), adminID, request, c.IP())
	if err != nil {
		return h.handleError(c, err)
	}

	h.setAuthCookie(c, session.Token)

	return utils.Success(c, fiber.StatusOK, "impersonation started successfully", fiber.Map{
		"user": session.User,
		"impersonation": fiber.Map{
			"active":          true,
			"impersonator_id": adminID.String(),
			"target_user_id":  session.User.ID,
		},
	})
}

func (h *AdminHandler) StopImpersonation(c *fiber.Ctx) error {
	targetUserID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	impersonatorID, err := middleware.GetImpersonatorID(c)
	if err != nil {
		return h.handleError(c, err)
	}

	session, err := h.impersonationService.Stop(c.UserContext(), targetUserID, impersonatorID, c.IP())
	if err != nil {
		return h.handleError(c, err)
	}

	h.setAuthCookie(c, session.Token)

	return utils.Success(c, fiber.StatusOK, "impersonation ended successfully", fiber.Map{
		"user": session.User,
		"impersonation": fiber.Map{
			"active": false,
		},
	})
}

func (h *AdminHandler) handleError(c *fiber.Ctx, err error) error {
	var validationErr *models.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return utils.Error(c, fiber.StatusBadRequest, validationErr.Error(), validationErr.Fields)
	case errors.Is(err, models.ErrAdminOnly), errors.Is(err, models.ErrImpersonationNotActive):
		return utils.Error(c, fiber.StatusForbidden, err.Error(), nil)
	case errors.Is(err, models.ErrInvalidImpersonationPassword):
		return utils.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
	case errors.Is(err, models.ErrUserNotFound):
		return utils.Error(c, fiber.StatusNotFound, "user not found", nil)
	case errors.Is(err, models.ErrCannotImpersonateSelf), errors.Is(err, models.ErrImpersonationNotConfigured):
		return utils.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	default:
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}
}

func (h *AdminHandler) setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     h.jwtConfig.CookieName,
		Value:    token,
		HTTPOnly: true,
		Secure:   h.jwtConfig.SecureCookie,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
		MaxAge:   int(h.jwtConfig.Expiration.Seconds()),
		Expires:  time.Now().Add(h.jwtConfig.Expiration),
	})
}

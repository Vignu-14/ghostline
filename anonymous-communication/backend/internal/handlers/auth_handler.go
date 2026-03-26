package handlers

import (
	"errors"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
	jwtConfig   config.JWTConfig
}

func NewAuthHandler(authService *services.AuthService, jwtConfig config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtConfig:   jwtConfig,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request models.RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	session, err := h.authService.Register(c.UserContext(), request)
	if err != nil {
		return h.handleError(c, err)
	}

	h.setAuthCookie(c, session.Token)

	return utils.Success(c, fiber.StatusCreated, "account created successfully", fiber.Map{
		"user": session.User,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request models.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	session, err := h.authService.Login(c.UserContext(), request, c.IP(), c.Get(fiber.HeaderUserAgent))
	if err != nil {
		return h.handleError(c, err)
	}

	h.setAuthCookie(c, session.Token)

	return utils.Success(c, fiber.StatusOK, "login successful", fiber.Map{
		"user": session.User,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.clearAuthCookie(c)

	return utils.Success(c, fiber.StatusOK, "logout successful", nil)
}

func (h *AuthHandler) handleError(c *fiber.Ctx, err error) error {
	var validationErr *models.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return utils.Error(c, fiber.StatusBadRequest, validationErr.Error(), validationErr.Fields)
	case errors.Is(err, models.ErrUsernameTaken), errors.Is(err, models.ErrEmailTaken):
		return utils.Error(c, fiber.StatusConflict, err.Error(), nil)
	case errors.Is(err, models.ErrInvalidCredentials):
		return utils.Error(c, fiber.StatusUnauthorized, "invalid credentials", nil)
	default:
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}
}

func (h *AuthHandler) setAuthCookie(c *fiber.Ctx, token string) {
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

func (h *AuthHandler) clearAuthCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     h.jwtConfig.CookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   h.jwtConfig.SecureCookie,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

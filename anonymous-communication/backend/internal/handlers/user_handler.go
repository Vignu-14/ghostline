package handlers

import (
	"errors"

	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	user, err := h.userService.GetByID(c.UserContext(), userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "user not found", nil)
		}

		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return utils.Success(c, fiber.StatusOK, "current user fetched successfully", fiber.Map{
		"user": user,
	})
}

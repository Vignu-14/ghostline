package handlers

import (
	"errors"

	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (h *UserHandler) Search(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		userID = uuid.Nil
	}

	query := c.Query("q")
	limit := c.QueryInt("limit", 8)
	if limit < 1 {
		limit = 8
	}
	if limit > 20 {
		limit = 20
	}

	users, err := h.userService.SearchByUsername(c.UserContext(), userID, query, limit)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return utils.Success(c, fiber.StatusOK, "users fetched successfully", fiber.Map{
		"users": users,
		"limit": limit,
		"query": query,
	})
}

func (h *UserHandler) Profile(c *fiber.Ctx) error {
	username := c.Params("username")
	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	profile, posts, err := h.userService.GetProfileByUsername(c.UserContext(), username, limit, offset)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return utils.Error(c, fiber.StatusNotFound, "user not found", nil)
		}

		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return utils.Success(c, fiber.StatusOK, "profile fetched successfully", fiber.Map{
		"profile": profile,
		"posts":   posts,
		"page":    page,
		"limit":   limit,
	})
}

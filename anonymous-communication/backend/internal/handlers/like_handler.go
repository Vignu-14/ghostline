package handlers

import (
	"errors"

	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type LikeHandler struct {
	likeService *services.LikeService
}

func NewLikeHandler(likeService *services.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

func (h *LikeHandler) Like(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	postID, err := utils.ParseUUID(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid post id", nil)
	}

	if err := h.likeService.Like(c.UserContext(), userID, postID); err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "post liked successfully", nil)
}

func (h *LikeHandler) Unlike(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	postID, err := utils.ParseUUID(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid post id", nil)
	}

	if err := h.likeService.Unlike(c.UserContext(), userID, postID); err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "post unliked successfully", nil)
}

func (h *LikeHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrPostNotFound):
		return utils.Error(c, fiber.StatusNotFound, "post not found", nil)
	case errors.Is(err, models.ErrCannotLikeOwnPost):
		return utils.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	default:
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}
}

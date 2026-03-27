package handlers

import (
	"errors"
	"mime/multipart"

	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/models"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) List(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	posts, err := h.postService.ListFeed(c.UserContext(), limit, offset)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return utils.Success(c, fiber.StatusOK, "posts fetched successfully", fiber.Map{
		"posts": posts,
		"page":  page,
		"limit": limit,
	})
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	var (
		file       multipart.File
		fileHeader *multipart.FileHeader
	)

	fileHeader, err = c.FormFile("image")
	if err == nil && fileHeader != nil {
		openedFile, openErr := fileHeader.Open()
		if openErr != nil {
			return utils.Error(c, fiber.StatusBadRequest, "unable to read uploaded file", nil)
		}
		defer openedFile.Close()
		file = openedFile
	}

	post, err := h.postService.Create(c.UserContext(), userID, file, fileHeader, c.FormValue("caption"))
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusCreated, "post created successfully", fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) CreateUploadURL(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	var request models.CreatePostUploadRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	upload, err := h.postService.CreateUploadURL(c.UserContext(), userID, request)
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "upload url created successfully", fiber.Map{
		"upload": upload,
	})
}

func (h *PostHandler) CreateFromUploadedObject(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	var request models.CreatePostFromUploadRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	post, err := h.postService.CreateFromUploadedObject(c.UserContext(), userID, request)
	if err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusCreated, "post created successfully", fiber.Map{
		"post": post,
	})
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "authentication required", nil)
	}

	postID, err := utils.ParseUUID(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "invalid post id", nil)
	}

	if err := h.postService.Delete(c.UserContext(), userID, postID); err != nil {
		return h.handleError(c, err)
	}

	return utils.Success(c, fiber.StatusOK, "post deleted successfully", nil)
}

func (h *PostHandler) handleError(c *fiber.Ctx, err error) error {
	var validationErr *models.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return utils.Error(c, fiber.StatusBadRequest, validationErr.Error(), validationErr.Fields)
	case errors.Is(err, models.ErrStorageNotConfigured):
		return utils.Error(c, fiber.StatusServiceUnavailable, "storage is not configured", nil)
	case errors.Is(err, models.ErrPostNotFound):
		return utils.Error(c, fiber.StatusNotFound, "post not found", nil)
	case errors.Is(err, models.ErrForbidden):
		return utils.Error(c, fiber.StatusForbidden, "forbidden", nil)
	default:
		return utils.Error(c, fiber.StatusInternalServerError, "internal server error", nil)
	}
}

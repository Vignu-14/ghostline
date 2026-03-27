package middleware

import (
	"errors"

	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError
	message := "internal server error"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		statusCode = fiberErr.Code
		if statusCode < fiber.StatusInternalServerError {
			message = fiberErr.Message
		}
	}

	requestID := GetRequestID(c)
	if requestID != "" {
		c.Set(RequestIDHeader, requestID)
	}

	return c.Status(statusCode).JSON(utils.APIResponse{
		Status: "error",
		Error:  message,
	})
}

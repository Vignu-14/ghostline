package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    any               `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string, details map[string]string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "error",
		Error:   message,
		Details: details,
	})
}

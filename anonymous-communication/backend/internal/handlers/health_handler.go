package handlers

import (
	"context"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ok",
		"service":   config.AppName,
		"timestamp": time.Now().UTC(),
	})
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbHealth, err := database.Health(ctx, h.db)
	statusCode := fiber.StatusOK
	status := "ok"
	if err != nil {
		statusCode = fiber.StatusServiceUnavailable
		status = "degraded"
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":    status,
		"service":   config.AppName,
		"checks":    fiber.Map{"database": dbHealth},
		"timestamp": time.Now().UTC(),
	})
}

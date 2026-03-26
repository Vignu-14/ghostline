package middleware

import (
	"anonymous-communication/backend/internal/config"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCORS(cfg config.CORSConfig) fiber.Handler {
	return fibercors.New(fibercors.Config{
		AllowOrigins:     cfg.AllowedOrigin,
		AllowCredentials: cfg.AllowCredentials,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
	})
}

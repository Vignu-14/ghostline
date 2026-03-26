package routes

import (
	"anonymous-communication/backend/internal/handlers"
	"anonymous-communication/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(
	app *fiber.App,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	app.Get("/health", healthHandler.Live)

	api := app.Group("/api")
	registerAPIRoutes(api, healthHandler, authHandler, userHandler, jwtMiddleware)
	registerWebSocketRoutes(api)
}

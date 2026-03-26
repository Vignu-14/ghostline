package routes

import (
	"anonymous-communication/backend/internal/handlers"
	"anonymous-communication/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func registerAPIRoutes(
	api fiber.Router,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	api.Get("/health", healthHandler.Live)
	api.Get("/health/ready", healthHandler.Ready)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", jwtMiddleware.RequireAuth, userHandler.Me)

	users := api.Group("/users")
	users.Get("/me", jwtMiddleware.RequireAuth, userHandler.Me)
}

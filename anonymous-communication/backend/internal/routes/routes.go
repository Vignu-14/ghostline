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
	adminHandler *handlers.AdminHandler,
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
	likeHandler *handlers.LikeHandler,
	chatHandler *handlers.ChatHandler,
	websocketHandler *handlers.WebSocketHandler,
	jwtMiddleware *middleware.JWTMiddleware,
	adminMiddleware *middleware.AdminMiddleware,
) {
	app.Get("/health", healthHandler.Live)

	api := app.Group("/api")
	registerAPIRoutes(api, healthHandler, authHandler, adminHandler, userHandler, postHandler, likeHandler, chatHandler, jwtMiddleware, adminMiddleware)
	registerWebSocketRoutes(app, websocketHandler)
}

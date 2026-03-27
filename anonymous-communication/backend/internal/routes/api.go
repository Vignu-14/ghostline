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
	adminHandler *handlers.AdminHandler,
	userHandler *handlers.UserHandler,
	postHandler *handlers.PostHandler,
	likeHandler *handlers.LikeHandler,
	chatHandler *handlers.ChatHandler,
	jwtMiddleware *middleware.JWTMiddleware,
	adminMiddleware *middleware.AdminMiddleware,
) {
	api.Get("/health", healthHandler.Live)
	api.Get("/health/ready", healthHandler.Ready)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", jwtMiddleware.RequireAuth, userHandler.Me)

	api.Get("/posts", postHandler.List)

	users := api.Group("/users")
	users.Get("/me", jwtMiddleware.RequireAuth, userHandler.Me)

	posts := api.Group("/posts")
	posts.Post("/", jwtMiddleware.RequireAuth, postHandler.Create)
	posts.Delete("/:id", jwtMiddleware.RequireAuth, postHandler.Delete)
	posts.Post("/:id/like", jwtMiddleware.RequireAuth, likeHandler.Like)
	posts.Delete("/:id/like", jwtMiddleware.RequireAuth, likeHandler.Unlike)

	messages := api.Group("/messages", jwtMiddleware.RequireAuth)
	messages.Get("/conversations", chatHandler.ListConversations)
	messages.Get("/:userId", chatHandler.GetConversation)
	messages.Post("/", chatHandler.SendMessage)

	admin := api.Group("/admin", jwtMiddleware.RequireAuth)
	admin.Post("/impersonate", adminMiddleware.RequireAdmin, adminHandler.Impersonate)
	admin.Post("/impersonate/stop", adminHandler.StopImpersonation)
}

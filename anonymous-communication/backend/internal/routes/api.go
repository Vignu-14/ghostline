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
	rateLimiter *middleware.RateLimiter,
) {
	api.Get("/health", healthHandler.Live)
	api.Get("/health/ready", healthHandler.Ready)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", rateLimiter.Login(), authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", jwtMiddleware.RequireAuth, userHandler.Me)

	api.Get("/posts", postHandler.List)

	users := api.Group("/users")
	users.Get("/me", jwtMiddleware.RequireAuth, userHandler.Me)
	users.Get("/profile/:username", userHandler.Profile)
	users.Get("/search", userHandler.Search)

	posts := api.Group("/posts")
	posts.Post("/upload-url", jwtMiddleware.RequireAuth, rateLimiter.Uploads(), postHandler.CreateUploadURL)
	posts.Post("/finalize", jwtMiddleware.RequireAuth, postHandler.CreateFromUploadedObject)
	posts.Post("/", jwtMiddleware.RequireAuth, rateLimiter.Uploads(), postHandler.Create)
	posts.Delete("/:id", jwtMiddleware.RequireAuth, postHandler.Delete)
	posts.Post("/:id/like", jwtMiddleware.RequireAuth, rateLimiter.Likes(), likeHandler.Like)
	posts.Delete("/:id/like", jwtMiddleware.RequireAuth, rateLimiter.Likes(), likeHandler.Unlike)

	messages := api.Group("/messages", jwtMiddleware.RequireAuth)
	messages.Get("/conversations", chatHandler.ListConversations)
	messages.Post("/delete", chatHandler.DeleteMessages)
	messages.Post("/:userId/clear", chatHandler.ClearConversation)
	messages.Get("/:userId", chatHandler.GetConversation)
	messages.Post("/", rateLimiter.Messages(), chatHandler.SendMessage)

	admin := api.Group("/admin", jwtMiddleware.RequireAuth)
	admin.Post("/impersonate", adminMiddleware.RequireAdmin, adminHandler.Impersonate)
	admin.Post("/impersonate/stop", adminHandler.StopImpersonation)
}

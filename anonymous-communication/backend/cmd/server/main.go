package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/database"
	"anonymous-communication/backend/internal/handlers"
	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/repositories"
	"anonymous-communication/backend/internal/routes"
	"anonymous-communication/backend/internal/services"
	"anonymous-communication/backend/internal/websocket"

	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("backend exited", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := loadEnvFile(); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dbPool, err := database.Connect(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer database.Close(dbPool)

	app := fiber.New(fiber.Config{
		AppName:      config.AppName,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(middleware.RequestID)
	app.Use(middleware.RequestLogger)
	app.Use(fiberrecover.New())
	app.Use(middleware.NewSecureHeaders(cfg.Server.Environment))
	app.Use(middleware.NewCORS(cfg.CORS))

	userRepository := repositories.NewUserRepository(dbPool)
	authLogRepository := repositories.NewAuthLogRepository(dbPool)
	adminRepository := repositories.NewAdminRepository(dbPool)
	postRepository := repositories.NewPostRepository(dbPool)
	likeRepository := repositories.NewLikeRepository(dbPool)
	messageRepository := repositories.NewMessageRepository(dbPool)

	authService := services.NewAuthService(userRepository, authLogRepository, cfg.JWT)
	userService := services.NewUserService(userRepository, postRepository)
	impersonationService := services.NewImpersonationService(userRepository, adminRepository, cfg.JWT)
	uploadService := services.NewUploadService(cfg.Storage)
	postService := services.NewPostService(postRepository, uploadService)
	likeService := services.NewLikeService(postRepository, likeRepository)
	chatService := services.NewChatService(messageRepository, userRepository)
	websocketHub := websocket.NewHub()
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)

	healthHandler := handlers.NewHealthHandler(dbPool)
	authHandler := handlers.NewAuthHandler(authService, cfg.JWT)
	adminHandler := handlers.NewAdminHandler(impersonationService, cfg.JWT)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	likeHandler := handlers.NewLikeHandler(likeService)
	chatHandler := handlers.NewChatHandler(chatService)
	websocketHandler := handlers.NewWebSocketHandler(chatService, websocketHub, cfg.JWT, rateLimiter)
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT)
	adminMiddleware := middleware.NewAdminMiddleware()

	routes.Register(app, healthHandler, authHandler, adminHandler, userHandler, postHandler, likeHandler, chatHandler, websocketHandler, jwtMiddleware, adminMiddleware, rateLimiter)

	serverErrors := make(chan error, 1)
	go func() {
		address := ":" + cfg.Server.Port
		slog.Info("starting backend server", "address", address, "environment", cfg.Server.Environment)
		serverErrors <- app.Listen(address)
	}()

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	case <-shutdownSignal.Done():
		slog.Info("shutting down backend server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown fiber app: %w", err)
		}

		return nil
	}
}

func loadEnvFile() error {
	candidates := []string{
		".env",
		filepath.Join("backend", ".env"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return fmt.Errorf("stat %s: %w", path, err)
		}

		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}
	}

	return nil
}

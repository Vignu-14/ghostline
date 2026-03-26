package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/database"
	"anonymous-communication/backend/internal/handlers"
	"anonymous-communication/backend/internal/middleware"
	"anonymous-communication/backend/internal/repositories"
	"anonymous-communication/backend/internal/routes"
	"anonymous-communication/backend/internal/services"

	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if err := run(); err != nil {
		slog.Error("backend exited", "error", err)
		os.Exit(1)
	}
}

func run() error {
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
	})

	app.Use(fiberrecover.New())
	app.Use(middleware.NewCORS(cfg.CORS))

	userRepository := repositories.NewUserRepository(dbPool)
	authLogRepository := repositories.NewAuthLogRepository(dbPool)

	authService := services.NewAuthService(userRepository, authLogRepository, cfg.JWT)
	userService := services.NewUserService(userRepository)

	healthHandler := handlers.NewHealthHandler(dbPool)
	authHandler := handlers.NewAuthHandler(authService, cfg.JWT)
	userHandler := handlers.NewUserHandler(userService)
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWT)

	routes.Register(app, healthHandler, authHandler, userHandler, jwtMiddleware)

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

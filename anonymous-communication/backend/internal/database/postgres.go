package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"anonymous-communication/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer pingCancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}


func Health(ctx context.Context, pool *pgxpool.Pool) (map[string]any, error) {
	if pool == nil {
		return map[string]any{
			"status": "down",
		}, errors.New("database pool is not initialized")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return map[string]any{
			"status": "down",
			"error":  err.Error(),
		}, fmt.Errorf("ping postgres: %w", err)
	}

	stats := pool.Stat()

	return map[string]any{
		"status":         "up",
		"max_conns":      stats.MaxConns(),
		"acquired_conns": stats.AcquiredConns(),
		"idle_conns":     stats.IdleConns(),
		"total_conns":    stats.TotalConns(),
	}, nil
}

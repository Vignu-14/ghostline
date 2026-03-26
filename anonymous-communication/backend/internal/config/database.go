package config

import (
	"fmt"
	"strings"
	"time"
)

type DatabaseConfig struct {
	URL               string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

func (c DatabaseConfig) Validate() error {
	if strings.TrimSpace(c.URL) == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.MaxConns < 1 {
		return fmt.Errorf("DB_MAX_CONNECTIONS must be greater than 0")
	}

	if c.MinConns < 0 {
		return fmt.Errorf("DB_MIN_CONNECTIONS cannot be negative")
	}

	if c.MinConns > c.MaxConns {
		return fmt.Errorf("DB_MIN_CONNECTIONS cannot be greater than DB_MAX_CONNECTIONS")
	}

	return nil
}

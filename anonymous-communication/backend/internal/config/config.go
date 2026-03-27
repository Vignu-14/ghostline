package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	Storage   StorageConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port            string
	Environment     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type JWTConfig struct {
	Secret       string
	Expiration   time.Duration
	CookieName   string
	SecureCookie bool
}

type CORSConfig struct {
	AllowedOrigin    string
	AllowCredentials bool
}

type RateLimitConfig struct {
	LoginAttempts int
	LoginWindow   time.Duration
	UploadCount   int
	UploadWindow  time.Duration
	MessageCount  int
	MessageWindow time.Duration
	LikeCount     int
	LikeWindow    time.Duration
}

func Load() (*Config, error) {
	environment := getEnv("ENVIRONMENT", DefaultEnvironment)

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", DefaultPort),
			Environment:     environment,
			ReadTimeout:     getEnvDurationSeconds("READ_TIMEOUT_SECONDS", DefaultReadTimeout),
			WriteTimeout:    getEnvDurationSeconds("WRITE_TIMEOUT_SECONDS", DefaultWriteTimeout),
			IdleTimeout:     getEnvDurationSeconds("IDLE_TIMEOUT_SECONDS", DefaultIdleTimeout),
			ShutdownTimeout: getEnvDurationSeconds("SHUTDOWN_TIMEOUT_SECONDS", DefaultShutdownTimeout),
		},
		Database: DatabaseConfig{
			URL:               strings.TrimSpace(os.Getenv("DATABASE_URL")),
			MaxConns:          int32(getEnvInt("DB_MAX_CONNECTIONS", int(DefaultDBMaxConns))),
			MinConns:          int32(getEnvInt("DB_MIN_CONNECTIONS", int(DefaultDBMinConns))),
			MaxConnLifetime:   getEnvDurationMinutes("DB_MAX_CONN_LIFETIME_MINUTES", DefaultDBMaxConnLifetime),
			MaxConnIdleTime:   getEnvDurationMinutes("DB_MAX_CONN_IDLE_MINUTES", DefaultDBMaxConnIdleTime),
			HealthCheckPeriod: getEnvDurationSeconds("DB_HEALTH_CHECK_SECONDS", DefaultDBHealthCheckPeriod),
			ConnectTimeout:    getEnvDurationSeconds("DB_CONNECT_TIMEOUT_SECONDS", DefaultDBConnectTimeout),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", defaultJWTSecret(environment)),
			Expiration:   getEnvDurationMinutes("JWT_EXPIRATION_MINUTES", DefaultJWTExpiration),
			CookieName:   getEnv("AUTH_COOKIE_NAME", DefaultCookieName),
			SecureCookie: getEnvBool("COOKIE_SECURE", isProduction(environment)),
		},
		CORS: CORSConfig{
			AllowedOrigin:    getEnv("ALLOWED_ORIGIN", DefaultAllowedOrigin),
			AllowCredentials: true,
		},
		Storage: StorageConfig{
			SupabaseURL:        strings.TrimSpace(os.Getenv("SUPABASE_URL")),
			SupabaseServiceKey: strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_KEY")),
			BucketName:         getEnv("STORAGE_BUCKET_NAME", DefaultStorageBucketName),
		},
		RateLimit: RateLimitConfig{
			LoginAttempts: getEnvInt("RATE_LIMIT_LOGIN_ATTEMPTS", DefaultLoginRateLimitAttempts),
			LoginWindow:   getEnvDurationMinutes("RATE_LIMIT_LOGIN_WINDOW_MINUTES", DefaultLoginRateLimitWindow),
			UploadCount:   getEnvInt("RATE_LIMIT_UPLOAD_COUNT", DefaultUploadRateLimitCount),
			UploadWindow:  getEnvDurationMinutes("RATE_LIMIT_UPLOAD_WINDOW_MINUTES", DefaultUploadRateLimitWindow),
			MessageCount:  getEnvInt("RATE_LIMIT_MESSAGE_COUNT", DefaultMessageRateLimitCount),
			MessageWindow: getEnvDurationSeconds("RATE_LIMIT_MESSAGE_WINDOW_SECONDS", DefaultMessageRateLimitWindow),
			LikeCount:     getEnvInt("RATE_LIMIT_LIKE_COUNT", DefaultLikeRateLimitCount),
			LikeWindow:    getEnvDurationMinutes("RATE_LIMIT_LIKE_WINDOW_MINUTES", DefaultLikeRateLimitWindow),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Server.Port) == "" {
		return fmt.Errorf("PORT is required")
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	if strings.TrimSpace(c.JWT.CookieName) == "" {
		return fmt.Errorf("AUTH_COOKIE_NAME is required")
	}

	if isProduction(c.Server.Environment) && strings.TrimSpace(c.JWT.Secret) == "" {
		return fmt.Errorf("JWT_SECRET is required in production")
	}

	if err := c.RateLimit.Validate(); err != nil {
		return err
	}

	return nil
}

func (c RateLimitConfig) Validate() error {
	switch {
	case c.LoginAttempts < 1:
		return fmt.Errorf("RATE_LIMIT_LOGIN_ATTEMPTS must be greater than 0")
	case c.LoginWindow <= 0:
		return fmt.Errorf("RATE_LIMIT_LOGIN_WINDOW_MINUTES must be greater than 0")
	case c.UploadCount < 1:
		return fmt.Errorf("RATE_LIMIT_UPLOAD_COUNT must be greater than 0")
	case c.UploadWindow <= 0:
		return fmt.Errorf("RATE_LIMIT_UPLOAD_WINDOW_MINUTES must be greater than 0")
	case c.MessageCount < 1:
		return fmt.Errorf("RATE_LIMIT_MESSAGE_COUNT must be greater than 0")
	case c.MessageWindow <= 0:
		return fmt.Errorf("RATE_LIMIT_MESSAGE_WINDOW_SECONDS must be greater than 0")
	case c.LikeCount < 1:
		return fmt.Errorf("RATE_LIMIT_LIKE_COUNT must be greater than 0")
	case c.LikeWindow <= 0:
		return fmt.Errorf("RATE_LIMIT_LIKE_WINDOW_MINUTES must be greater than 0")
	default:
		return nil
	}
}

func defaultJWTSecret(environment string) string {
	if isProduction(environment) {
		return ""
	}

	return DefaultJWTSecret
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDurationSeconds(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}

	return time.Duration(parsed) * time.Second
}

func getEnvDurationMinutes(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return fallback
	}

	return time.Duration(parsed) * time.Minute
}

func isProduction(environment string) bool {
	return strings.EqualFold(strings.TrimSpace(environment), "production")
}

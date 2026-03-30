package config

import "time"

const (
	AppName                  = "anonymous-communication-backend"
	DefaultPort              = "3000"
	DefaultEnvironment       = "development"
	DefaultAllowedOrigin     = "http://localhost:5173"
	DefaultJWTSecret         = "change-me-for-local-development"
	DefaultStorageBucketName = "user-uploads"
	DefaultCookieName        = "auth_token"

	DefaultReadTimeout     = 15 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
	DefaultJWTExpiration   = 15 * time.Minute

	DefaultDBConnectTimeout          = 5 * time.Second
	DefaultDBMaxConnLifetime         = 1 * time.Hour
	DefaultDBMaxConnIdleTime         = 15 * time.Minute
	DefaultDBHealthCheckPeriod       = 30 * time.Second
	DefaultDBMaxConns          int32 = 25
	DefaultDBMinConns          int32 = 5

	DefaultLoginRateLimitAttempts    = 5
	DefaultLoginRateLimitWindow      = 15 * time.Minute
	DefaultRegisterRateLimitAttempts = 3
	DefaultRegisterRateLimitWindow   = 60 * time.Minute
	DefaultUploadRateLimitCount      = 10
	DefaultUploadRateLimitWindow     = 1 * time.Hour
	DefaultMessageRateLimitCount     = 10
	DefaultMessageRateLimitWindow    = 1 * time.Second
	DefaultLikeRateLimitCount        = 100
	DefaultLikeRateLimitWindow       = 1 * time.Hour
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

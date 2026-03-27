package middleware

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"anonymous-communication/backend/internal/config"
	"anonymous-communication/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	config      config.RateLimitConfig
	mu          sync.Mutex
	entries     map[string]rateLimitEntry
	lastCleanup time.Time
}

type rateLimitEntry struct {
	Count   int
	ResetAt time.Time
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:      cfg,
		entries:     make(map[string]rateLimitEntry),
		lastCleanup: time.Now(),
	}
}

func (r *RateLimiter) Login() fiber.Handler {
	return r.limitByIP(
		"login",
		r.config.LoginAttempts,
		r.config.LoginWindow,
		"too many login attempts. try again in 15 minutes.",
	)
}

func (r *RateLimiter) Uploads() fiber.Handler {
	return r.limitByActor(
		"uploads",
		r.config.UploadCount,
		r.config.UploadWindow,
		"upload rate limit exceeded. try again later.",
	)
}

func (r *RateLimiter) Messages() fiber.Handler {
	return r.limitByActor(
		"messages",
		r.config.MessageCount,
		r.config.MessageWindow,
		"message rate limit exceeded. slow down and try again.",
	)
}

func (r *RateLimiter) Likes() fiber.Handler {
	return r.limitByActor(
		"likes",
		r.config.LikeCount,
		r.config.LikeWindow,
		"like rate limit exceeded. try again later.",
	)
}

func (r *RateLimiter) AllowMessageForUser(userID string) (bool, time.Duration) {
	return r.allow("messages", "user:"+strings.TrimSpace(userID), r.config.MessageCount, r.config.MessageWindow)
}

func (r *RateLimiter) limitByIP(bucket string, limit int, window time.Duration, message string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return r.enforce(c, bucket, "ip:"+c.IP(), limit, window, message)
	}
}

func (r *RateLimiter) limitByActor(bucket string, limit int, window time.Duration, message string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return r.enforce(c, bucket, actorKey(c), limit, window, message)
	}
}

func (r *RateLimiter) enforce(c *fiber.Ctx, bucket, key string, limit int, window time.Duration, message string) error {
	allowed, retryAfter := r.allow(bucket, key, limit, window)
	if allowed {
		return c.Next()
	}

	retryAfterSeconds := int(math.Ceil(retryAfter.Seconds()))
	if retryAfterSeconds < 1 {
		retryAfterSeconds = 1
	}
	c.Set("Retry-After", strconv.Itoa(retryAfterSeconds))

	return utils.Error(c, fiber.StatusTooManyRequests, message, nil)
}

func (r *RateLimiter) allow(bucket, key string, limit int, window time.Duration) (bool, time.Duration) {
	now := time.Now()
	cacheKey := bucket + ":" + key

	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanupExpiredLocked(now)

	entry, exists := r.entries[cacheKey]
	if !exists || !now.Before(entry.ResetAt) {
		r.entries[cacheKey] = rateLimitEntry{
			Count:   1,
			ResetAt: now.Add(window),
		}
		return true, window
	}

	if entry.Count >= limit {
		return false, time.Until(entry.ResetAt)
	}

	entry.Count++
	r.entries[cacheKey] = entry
	return true, time.Until(entry.ResetAt)
}

func (r *RateLimiter) cleanupExpiredLocked(now time.Time) {
	if now.Sub(r.lastCleanup) < time.Minute {
		return
	}

	for key, entry := range r.entries {
		if !now.Before(entry.ResetAt) {
			delete(r.entries, key)
		}
	}

	r.lastCleanup = now
}

func actorKey(c *fiber.Ctx) string {
	if claims, ok := GetClaims(c); ok {
		userID := strings.TrimSpace(claims.UserID)
		if userID != "" {
			return "user:" + userID
		}
	}

	return "ip:" + c.IP()
}

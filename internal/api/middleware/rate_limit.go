package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"

	"social-media-app/internal/config"
	"social-media-app/internal/utils"
)

type RateLimiter struct {
	client *redis.Client
	cfg    *config.RateLimitConfig
	ctx    context.Context
}

func NewRateLimiter(cfg *config.Config) *RateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &RateLimiter{
		client: client,
		cfg:    &cfg.RateLimit,
		ctx:    context.Background(),
	}
}

// Rate limit based on authenticated user ID
func (rl *RateLimiter) RateLimitByUser(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		key := fmt.Sprintf("rate_limit:user:%d:%s", userID.(uint), action)

		if allowed, resetTime, err := rl.checkRateLimit(key); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Rate limiting error")
			c.Abort()
			return
		} else if !allowed {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rl.cfg.Requests))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
			utils.ErrorResponse(c, http.StatusTooManyRequests,
				fmt.Sprintf("Rate limit exceeded for %s. Try again later.", action))
			c.Abort()
			return
		}

		c.Next()
	}
}

// Rate limit based on IP address
func (rl *RateLimiter) RateLimitByIP(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:ip:%s:%s", clientIP, action)

		if allowed, resetTime, err := rl.checkRateLimit(key); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Rate limiting error")
			c.Abort()
			return
		} else if !allowed {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rl.cfg.Requests))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
			utils.ErrorResponse(c, http.StatusTooManyRequests,
				fmt.Sprintf("Rate limit exceeded for %s. Try again later.", action))
			c.Abort()
			return
		}

		c.Next()
	}
}

// Rate limit by user if authenticated, otherwise by IP
func (rl *RateLimiter) RateLimitByUserOrIP(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		if userID, exists := c.Get("user_id"); exists {
			key = fmt.Sprintf("rate_limit:user:%d:%s", userID.(uint), action)
		} else {
			clientIP := c.ClientIP()
			key = fmt.Sprintf("rate_limit:ip:%s:%s", clientIP, action)
		}

		if allowed, resetTime, err := rl.checkRateLimit(key); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Rate limiting error")
			c.Abort()
			return
		} else if !allowed {
			c.Header("X-RateLimit-Limit", strconv.Itoa(rl.cfg.Requests))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
			utils.ErrorResponse(c, http.StatusTooManyRequests,
				fmt.Sprintf("Rate limit exceeded for %s. Try again later.", action))
			c.Abort()
			return
		}
		c.Next()
	}
}

// Check rate limit using sliding window rate limiting using Redis
func (rl *RateLimiter) checkRateLimit(key string) (bool, int64, error) {
	now := time.Now().Unix()
	windowStart := now - int64(rl.cfg.Window.Seconds())

	pipe := rl.client.Pipeline()

	// Remove expired entries
	pipe.ZRemRangeByScore(rl.ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count current requests in window
	countCmd := pipe.ZCard(rl.ctx, key)

	// Add current request
	pipe.ZAdd(rl.ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Set expiry for the key
	pipe.Expire(rl.ctx, key, rl.cfg.Window)

	_, err := pipe.Exec(rl.ctx)
	if err != nil {
		return false, 0, err
	}

	currentCount := countCmd.Val()

	// Check if limit exceeded (subtract 1 because we already added current request)
	if currentCount > int64(rl.cfg.Requests) {
		// Remove the request we just added since it exceeds limit
		rl.client.ZRem(rl.ctx, key, fmt.Sprintf("%d", now))
		return false, now + int64(rl.cfg.Window.Seconds()), nil
	}

	return true, 0, nil
}

// Custom rate limiter with different limits
type CustomRateLimitConfig struct {
	Requests int
	Window   time.Duration
}

func (rl *RateLimiter) CustomRateLimit(action string, config CustomRateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		if userID, exists := c.Get("user_id"); exists {
			key = fmt.Sprintf("rate_limit:user:%d:%s", userID.(uint), action)
		} else {
			clientIP := c.ClientIP()
			key = fmt.Sprintf("rate_limit:ip:%s:%s", clientIP, action)
		}

		if allowed, resetTime, err := rl.checkCustomRateLimit(key, config); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Rate limiting error")
			c.Abort()
			return
		} else if !allowed {
			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
			utils.ErrorResponse(c, http.StatusTooManyRequests,
				fmt.Sprintf("Rate limit exceeded for %s. Try again later.", action))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) checkCustomRateLimit(key string, config CustomRateLimitConfig) (bool, int64, error) {
	now := time.Now().Unix()
	windowStart := now - int64(config.Window.Seconds())

	pipe := rl.client.Pipeline()

	// Remove expired entries
	pipe.ZRemRangeByScore(rl.ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count current requests in window
	countCmd := pipe.ZCard(rl.ctx, key)

	// Add current request
	pipe.ZAdd(rl.ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Set expiry for the key
	pipe.Expire(rl.ctx, key, config.Window)

	_, err := pipe.Exec(rl.ctx)
	if err != nil {
		return false, 0, err
	}

	currentCount := countCmd.Val()

	// Check if limit exceeded
	if currentCount > int64(config.Requests) {
		// Remove the request we just added since it exceeds limit
		rl.client.ZRem(rl.ctx, key, fmt.Sprintf("%d", now))
		return false, now + int64(config.Window.Seconds()), nil
	}

	return true, 0, nil
}

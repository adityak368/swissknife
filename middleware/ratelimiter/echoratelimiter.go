package ratelimiter

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RateLimitConfig defines the rate limiter config
type RateLimitConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
	Limit   float64
}

// DefaultRateLimitConfig defines the default limiter config
var DefaultRateLimitConfig = RateLimitConfig{
	Skipper: middleware.DefaultSkipper,
	Limit:   3,
}

// RateLimitMiddleware returns a middleware for echo with default config
func RateLimitMiddleware() echo.MiddlewareFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig)
}

// RateLimitWithConfig returns a middleware for echo with config
func RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimitConfig.Skipper
	}
	limiter := tollbooth.NewLimiter(config.Limit, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			httpError := tollbooth.LimitByRequest(limiter, c.Response(), c.Request())
			if httpError != nil {
				return echo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}

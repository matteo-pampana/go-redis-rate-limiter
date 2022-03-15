package ratelimiter

import "time"

const (
	DEFAULT_KEY              = "default"
	DEFAULT_MAX_REQUESTS     = 10
	DEFAULT_REFRESH_INTERVAL = 1 * time.Minute
)

// RateLimiter contains the rate limiter's configuration
type RateLimiterConfig struct {
	KeyItems        []string
	MaxRequests     int
	RefreshInterval time.Duration
}

// IsEmpty checks if the configuration is empty
func (cfg RateLimiterConfig) IsEmpty() bool {
	return len(cfg.KeyItems) == 0
}

// SetDefault sets the default configuration (bucket default, 10 reqs/minute)
func (cfg *RateLimiterConfig) SetDefault() {
	cfg.KeyItems = []string{DEFAULT_KEY}
	cfg.MaxRequests = DEFAULT_MAX_REQUESTS
	cfg.RefreshInterval = DEFAULT_REFRESH_INTERVAL
}

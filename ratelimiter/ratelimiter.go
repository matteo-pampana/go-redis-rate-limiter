package ratelimiter

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrTooManyRequests = errors.New("max requests limit reached")

	keySeparator = "#"
)

type store interface {
	GetCounter(ctx context.Context, requestKey string) (int, error)
	IncreaseWithTTL(ctx context.Context, requestKey string, ttl time.Duration) error
}

// RateLimiter represents a rate limiter
type RateLimiter struct {
	store  store
	config RateLimiterConfig
}

// NewRateLimiter creates a new rate limiter object
func NewRateLimiter(store store, config RateLimiterConfig) *RateLimiter {
	cfg := config
	if config.IsEmpty() {
		cfg.SetDefault()
	}
	return &RateLimiter{
		store:  store,
		config: cfg,
	}
}

// CheckRequest checks if the request is allowed,
// returns an error when reach the max number of requests
func (r RateLimiter) CheckRequest(ctx context.Context, keys []string) error {
	bucketKey := computeKeyFromItems(keys)
	val, err := r.store.GetCounter(ctx, bucketKey)
	if err != nil {
		return err
	}
	if val >= r.config.MaxRequests {
		return ErrTooManyRequests
	}

	err = r.store.IncreaseWithTTL(ctx, bucketKey, r.config.RefreshInterval)
	return err
}

func computeKeyFromItems(items []string) string {
	return strings.Join(items, keySeparator)
}

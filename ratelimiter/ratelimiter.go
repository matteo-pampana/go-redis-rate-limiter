package ratelimiter

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	errTooManyRequests = errors.New("max requests limit reached")

	keySeparator = "#"
)

type store interface {
	GetCounter(ctx context.Context, requestKey string) (int, error)
	IncreaseWithTTL(ctx context.Context, requestKey string, ttl time.Duration) error
}

// RateLimiter represents a rate limiter
type RateLimiter struct {
	store     store
	config    RateLimiterConfig
	bucketKey string
}

// NewRateLimiter creates a new rate limiter object
func NewRateLimiter(store store, config RateLimiterConfig) *RateLimiter {
	cfg := config
	if config.IsEmpty() {
		cfg.SetDefault()
	}
	bucketKey := computeKeyFromItems(cfg.KeyItems)
	return &RateLimiter{
		store:     store,
		config:    cfg,
		bucketKey: bucketKey,
	}
}

// CheckRequest checks if the request is allowed,
// returns an error when reach the max number of requests
func (r RateLimiter) CheckRequest(ctx context.Context) error {
	val, err := r.store.GetCounter(ctx, r.bucketKey)
	if err != nil {
		return err
	}
	if val >= r.config.MaxRequests {
		return errTooManyRequests
	}

	err = r.store.IncreaseWithTTL(ctx, r.bucketKey, r.config.RefreshInterval)
	return err
}

func computeKeyFromItems(items []string) string {
	return strings.Join(items, keySeparator)
}

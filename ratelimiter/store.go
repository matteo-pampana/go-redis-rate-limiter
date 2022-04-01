package ratelimiter

import (
	context "context"
	time "time"

	redis "github.com/go-redis/redis/v8"
)

type Store struct {
	redisClient *redis.Client
}

// NewStore creates a new store
func NewStore(redisClient *redis.Client) *Store {
	return &Store{
		redisClient: redisClient,
	}
}

// GetCounter returns the counter value for the given request key
func (s *Store) GetCounter(ctx context.Context, requestKey string) (int, error) {
	val, err := s.redisClient.Get(ctx, requestKey).Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}

// IncreaseWithTTL increases the counter value
// and set the TTL to the given bucket defined by the request key
func (s *Store) IncreaseWithTTL(ctx context.Context, requestKey string, ttl time.Duration) error {
	err := s.redisClient.Incr(ctx, requestKey).Err()
	if err != nil {
		return err
	}
	err = s.redisClient.Expire(ctx, requestKey, ttl).Err()
	return err
}

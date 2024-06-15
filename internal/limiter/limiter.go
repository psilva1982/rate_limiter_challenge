package limiter

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RateLimiter struct {
	client        *redis.Client
	IpRate        int
	TokenRate     int
	BlockDuration time.Duration
}

func NewRateLimiter(client *redis.Client, ipRate, tokenRate int, blockDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:        client,
		IpRate:        ipRate,
		TokenRate:     tokenRate,
		BlockDuration: blockDuration,
	}
}

func (r *RateLimiter) AllowRequest(identifier string, limit int) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limiter:%s", identifier)

	// Increment the counter
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// If it is the first request, set the counter expiration
	if count == 1 {
		r.client.Expire(ctx, key, time.Second)
	}

	// Check if the limit has been exceeded
	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}

func (r *RateLimiter) Block(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("block:%s", identifier)
	return r.client.Set(ctx, key, "blocked", r.BlockDuration).Err()
}

func (r *RateLimiter) IsBlocked(identifier string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("block:%s", identifier)
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "blocked", nil
}

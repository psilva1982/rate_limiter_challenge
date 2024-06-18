package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/psilva1982/rate_limiter_challenge/internal/limiter"
)

func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

type RedisRateLimiter struct {
	Client        *redis.Client
	IpRate        int
	TokenRate     int
	BlockDuration time.Duration
}

func NewRateLimiter() *RedisRateLimiter {
	client := InitRedis()
	ipRate, tokenRate, blockDuration := limiter.GetLimiterConfig()

	return &RedisRateLimiter{
		Client:        client,
		IpRate:        ipRate,
		TokenRate:     tokenRate,
		BlockDuration: blockDuration,
	}
}

func (r *RedisRateLimiter) AllowRequest(identifier string, limit int) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limiter:%s", identifier)

	// Increment the counter
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// If it is the first request, set the counter expiration
	if count == 1 {
		r.Client.Expire(ctx, key, time.Second)
	}

	// Check if the limit has been exceeded
	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}

func (r *RedisRateLimiter) Block(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("block:%s", identifier)
	return r.Client.Set(ctx, key, "blocked", r.BlockDuration).Err()
}

func (r *RedisRateLimiter) IsBlocked(identifier string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("block:%s", identifier)
	result, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "blocked", nil
}

func (r *RedisRateLimiter) GetIpRate() int {
	return r.IpRate
}

func (r *RedisRateLimiter) GetTokenRate() int {
	return r.TokenRate
}

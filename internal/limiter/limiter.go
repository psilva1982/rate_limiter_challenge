package limiter

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type RateLimiterClient interface {
	~*redis.Client | ~*redis.Client
}

type ChallangeLimiter interface {
	NewRateLimiter(client *redis.Client) *RateLimiter
	AllowRequest(identifier string, limit int) (bool, error)
}

type RateLimiter struct {
	client        *redis.Client
	IpRate        int
	TokenRate     int
	BlockDuration time.Duration
}

func GetLimiterConfig() (int, int, time.Duration) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ipRate, _ := strconv.Atoi(os.Getenv("IP_RATE_LIMIT"))
	tokenRate, _ := strconv.Atoi(os.Getenv("TOKEN_RATE_LIMIT"))
	blockDuration, _ := strconv.Atoi(os.Getenv("BLOCK_DURATION"))

	return ipRate, tokenRate, time.Duration(blockDuration) * time.Second
}

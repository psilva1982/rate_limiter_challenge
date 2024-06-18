package limiter

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type IRateLimiter interface {
	IsBlocked(identifier string) (bool, error)
	Block(identifier string) error
	AllowRequest(identifier string, limit int) (bool, error)
	GetIpRate() int
	GetTokenRate() int
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

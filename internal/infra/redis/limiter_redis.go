package redis

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

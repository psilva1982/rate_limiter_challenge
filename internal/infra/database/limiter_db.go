package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MySQLRateLimiter struct {
	db             *gorm.DB
	ipRateLimit    int
	tokenRateLimit int
	blockDuration  time.Duration
}

type RateLimiterEntry struct {
	ID        uint `gorm:"primaryKey"`
	Key       string
	Count     int64
	ExpiresAt time.Time
}

type TokenRateLimiterEntry struct {
	ID        uint `gorm:"primaryKey"`
	Token     string
	Count     int64
	ExpiresAt time.Time
}

func NewMySQLRateLimiter(db *gorm.DB, ipRateLimit, tokenRateLimit int, blockDuration time.Duration) *MySQLRateLimiter {
	return &MySQLRateLimiter{
		db:             db,
		ipRateLimit:    ipRateLimit,
		tokenRateLimit: tokenRateLimit,
		blockDuration:  blockDuration,
	}
}

func (rl *MySQLRateLimiter) AllowRequest(ip, token string) (bool, error) {
	ipKey := fmt.Sprintf("ip:%s", ip)
	tokenKey := fmt.Sprintf("token:%s", token)

	var ipEntry RateLimiterEntry
	result := rl.db.Where("key = ?", ipKey).FirstOrCreate(&ipEntry, RateLimiterEntry{Key: ipKey})
	if result.Error != nil {
		return false, result.Error
	}

	var tokenEntry TokenRateLimiterEntry
	result = rl.db.Where("token = ?", tokenKey).FirstOrCreate(&tokenEntry, TokenRateLimiterEntry{Token: tokenKey})
	if result.Error != nil {
		return false, result.Error
	}

	if ipEntry.Count >= int64(rl.ipRateLimit) || tokenEntry.Count >= int64(rl.tokenRateLimit) {
		// Block IP or token
		if ipEntry.Count >= int64(rl.ipRateLimit) {
			ipEntry.ExpiresAt = time.Now().Add(rl.blockDuration)
			rl.db.Save(&ipEntry)
		}
		if tokenEntry.Count >= int64(rl.tokenRateLimit) {
			tokenEntry.ExpiresAt = time.Now().Add(rl.blockDuration)
			rl.db.Save(&tokenEntry)
		}
		return false, nil
	}

	ipEntry.Count++
	tokenEntry.Count++

	rl.db.Save(&ipEntry)
	rl.db.Save(&tokenEntry)

	return true, nil
}

package database

import (
	"time"

	"github.com/psilva1982/rate_limiter_challenge/internal/limiter"
	"gorm.io/gorm"
)

type MySQLRateLimiter struct {
	DB            *gorm.DB
	IpRate        int
	TokenRate     int
	BlockDuration time.Duration
}

type RateLimit struct {
	ID         uint `gorm:"primaryKey"`
	Identifier string
	Count      int
	ExpiresAt  time.Time
}

type BlockList struct {
	ID           uint `gorm:"primaryKey"`
	Identifier   string
	BlockedUntil time.Time
}

func NewMySQLRateLimiter(db *gorm.DB) *MySQLRateLimiter {

	ipRate, tokenRate, blockDuration := limiter.GetLimiterConfig()

	return &MySQLRateLimiter{
		DB:            db,
		IpRate:        ipRate,
		TokenRate:     tokenRate,
		BlockDuration: blockDuration,
	}
}

func (r *MySQLRateLimiter) AllowRequest(identifier string, limit int) (bool, error) {
	var rateLimit RateLimit
	result := r.DB.Where("identifier = ?", identifier).First(&rateLimit)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	now := time.Now()
	if result.RowsAffected == 0 {
		rateLimit = RateLimit{
			Identifier: identifier,
			Count:      1,
			ExpiresAt:  now.Add(time.Second),
		}
		result = r.DB.Create(&rateLimit)
		if result.Error != nil {
			return false, result.Error
		}
	} else if now.After(rateLimit.ExpiresAt) {
		rateLimit.Count = 1
		rateLimit.ExpiresAt = now.Add(time.Second)
		result = r.DB.Save(&rateLimit)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		rateLimit.Count++
		if rateLimit.Count > limit {
			return false, nil
		}
		result = r.DB.Save(&rateLimit)
		if result.Error != nil {
			return false, result.Error
		}
	}

	return true, nil
}

func (r *MySQLRateLimiter) Block(identifier string) error {
	block := BlockList{
		Identifier:   identifier,
		BlockedUntil: time.Now().Add(r.BlockDuration),
	}
	result := r.DB.Create(&block)
	return result.Error
}

func (r *MySQLRateLimiter) IsBlocked(identifier string) (bool, error) {
	var block BlockList
	result := r.DB.Where("identifier = ?", identifier).First(&block)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return time.Now().Before(block.BlockedUntil), nil
}

func (r *MySQLRateLimiter) GetIpRate() int {
	return r.IpRate
}

func (r *MySQLRateLimiter) GetTokenRate() int {
	return r.TokenRate
}

package services

import (
	"math/rand"
	"time"

	"github.com/psilva1982/rate_limiter_challenge/internal/infra/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(email, password string) (*database.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &database.User{
		Email:    email,
		Password: string(hashedPassword),
		APIKey:   generateAPIKey(),
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*database.User, error) {
	var user database.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func generateAPIKey() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	apiKey := make([]byte, 32)
	for i := range apiKey {
		apiKey[i] = chars[rand.Intn(len(chars))]
	}
	return string(apiKey)
}

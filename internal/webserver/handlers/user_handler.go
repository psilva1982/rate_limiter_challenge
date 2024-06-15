package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/psilva1982/rate_limiter_challenge/internal/dto"
	"github.com/psilva1982/rate_limiter_challenge/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

// @Summary Create a new user
// @Description Create a new user with email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.UserInput true "User input"
// @Success 201
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Error creating user"
// @Router /signup [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput dto.UserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.CreateUser(userInput.Email, userInput.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// @Summary Get API key
// @Description Get API key for a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body dto.UserInput true "User input"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Invalid credentials"
// @Router /get-api-key [post]
func (h *UserHandler) GetAPIKey(w http.ResponseWriter, r *http.Request) {
	var userInput dto.UserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetUserByEmail(userInput.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"api_key": user.APIKey})
}

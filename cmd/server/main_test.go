package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/psilva1982/rate_limiter_challenge/internal/infra/database"
	"github.com/psilva1982/rate_limiter_challenge/internal/infra/redis"
	"github.com/psilva1982/rate_limiter_challenge/internal/limiter"
	"github.com/psilva1982/rate_limiter_challenge/internal/services"
	"github.com/psilva1982/rate_limiter_challenge/internal/webserver/handlers"
	customMiddleware "github.com/psilva1982/rate_limiter_challenge/internal/webserver/middleware"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
)

type Response struct {
	APIKey string `json:"api_key"`
}

func TestIPRateLimiter(t *testing.T) {

	ipRate, _, blockDur := limiter.GetLimiterConfig()
	rateLimiter := redis.NewRateLimiter()
	defaultHandler := handlers.NewDefaultHandler()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(customMiddleware.RateLimiterMiddleware(rateLimiter))
	r.Get("/", defaultHandler.BaseAccess)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Helper function to send a request, return the status code and the body
	sendRequest := func() (int, string) {
		req, _ := http.NewRequest("GET", ts.URL, nil)
		req.Header.Set("API_KEY", "") // no token, using only IP limitation
		resp, _ := http.DefaultClient.Do(req)
		bbody, _ := ioutil.ReadAll(resp.Body)
		body := string(bbody)
		return resp.StatusCode, body
	}

	var statusCode int
	var body string
	for i := 0; i < ipRate; i++ {
		statusCode, body = sendRequest()
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, "Request allowed", body)
	}

	// Send the next two requests that should be blocked
	for i := 0; i < 2; i++ {
		statusCode, body = sendRequest()
		assert.Equal(t, http.StatusTooManyRequests, statusCode)
	}
	// Wait the time required to reset the limit
	time.Sleep(blockDur)

	// Resend the first request after resetting the limit
	statusCode, body = sendRequest()
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Request allowed", body)
}

func TestTokenRateLimiter(t *testing.T) {

	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database")
	}
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)
	defaultHandler := handlers.NewDefaultHandler()

	_, tokenRate, blockDur := limiter.GetLimiterConfig()
	rateLimiter := redis.NewRateLimiter()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(customMiddleware.RateLimiterMiddleware(rateLimiter))
	r.Post("/signup", userHandler.CreateUser)
	r.Post("/get-api-key", userHandler.GetAPIKey)
	r.Get("/", defaultHandler.BaseAccess)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Fake User
	user := database.User{Email: "test@test.com", Password: "password", APIKey: "fake_key"}
	db.Create(&user)

	// Helper function to send a request, return the status code and the body
	sendRequest := func() (int, string) {
		req, _ := http.NewRequest("GET", ts.URL, nil)
		req.Header.Set("API_KEY", user.APIKey) // no token, using only IP limitation
		resp, _ := http.DefaultClient.Do(req)
		bbody, _ := ioutil.ReadAll(resp.Body)
		body := string(bbody)
		return resp.StatusCode, body
	}

	var statusCode int
	var body string
	for i := 0; i < tokenRate; i++ {
		statusCode, body = sendRequest()
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, "Request allowed", body)
	}

	// Send the next two requests that should be blocked
	for i := 0; i < 2; i++ {
		statusCode, body = sendRequest()
		assert.Equal(t, http.StatusTooManyRequests, statusCode)
	}
	// Wait the time required to reset the limit
	time.Sleep(blockDur)

	// Resend the first request after resetting the limit
	statusCode, body = sendRequest()
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "Request allowed", body)

	db.Delete(&user)
}

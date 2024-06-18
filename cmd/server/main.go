package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/psilva1982/rate_limiter_challenge/docs"
	"github.com/psilva1982/rate_limiter_challenge/internal/infra/database"
	"github.com/psilva1982/rate_limiter_challenge/internal/infra/redis"
	"github.com/psilva1982/rate_limiter_challenge/internal/services"
	"github.com/psilva1982/rate_limiter_challenge/internal/webserver/handlers"
	customMiddleware "github.com/psilva1982/rate_limiter_challenge/internal/webserver/middleware"
)

// @title Rate Limiter API
// @version 1.0
// @description This is a sample rate limiter server.
// @termsOfService http://swagger.io/terms/

// @contact.name Rate Limiter Challange
// @contact.url github.com/psilva1982/rate_limiter_challenge/
// @contact.email severos1982@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	godotenv.Load("/.env")

	// Startup MySQL
	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Redis RateLimiter
	rateLimiter := redis.NewRateLimiter()

	// Mysql RateLimiter
	//rateLimiter := database.NewMySQLRateLimiter(db)

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)
	defaultHandler := handlers.NewDefaultHandler()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(customMiddleware.RateLimiterMiddleware(rateLimiter))

	r.Post("/signup", userHandler.CreateUser)
	r.Post("/get-api-key", userHandler.GetAPIKey)
	r.Get("/", defaultHandler.BaseAccess)

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/doc.json"),
	))

	http.ListenAndServe(":8080", r)
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/psilva1982/rate_limiter_challenge/internal/infra/redis"
)

func RateLimiterMiddleware(rl *redis.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			token := r.Header.Get("API_KEY")
			identifier := ip
			limit := rl.IpRate

			if token != "" {
				identifier = token
				limit = rl.TokenRate
			}

			if isBlocked, _ := rl.IsBlocked(identifier); isBlocked {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			if allowed, _ := rl.AllowRequest(identifier, limit); !allowed {
				rl.Block(identifier)
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

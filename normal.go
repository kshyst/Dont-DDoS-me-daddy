package Daddy

import (
	"context"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"github.com/redis/go-redis/v9"
	"net"
	"net/http"
	"time"
)

func RateLimiter(next http.Handler, redisClient *redis.Client, options ...services.Option) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create the context with 10 second timeout
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// Check Redis availability first
		if err := redisClient.Ping(ctx).Err(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// Prepare request data
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write([]byte(err.Error()))
		}

		requestData := &models.ReqData{
			UserIp:         host,
			RequestAddress: r.RequestURI,
		}

		// Create the ratelimiter service
		service := services.NewService(db.CreateRedis(redisClient), options...)

		// Check rate limit
		if allowed := service.CheckAndStoreRate(ctx, requestData); !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limit exceeded"))
			return
		}

		// Goes to the next handler
		next.ServeHTTP(w, r)
	})
}

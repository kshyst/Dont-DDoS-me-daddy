package Daddy

import (
	"context"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"github.com/labstack/echo"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

func EchoMiddleware(redisClient *redis.Client, options ...services.Option) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a child context with timeout that inherits from the request context
			ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
			defer cancel()

			// Check Redis availability first
			if err := redisClient.Ping(ctx).Err(); err != nil {
				c.Response().Status = http.StatusServiceUnavailable
				log.Fatal("redis ping failed: ", err)
				return echo.ErrServiceUnavailable
			}

			// Prepare request data
			requestData := &models.ReqData{
				UserIp:         c.RealIP(),
				RequestAddress: c.Request().RequestURI,
			}

			// Create service instance
			service := services.NewService(db.CreateRedis(redisClient), options...)

			// Check rate limit
			if allowed := service.CheckAndStoreRate(ctx, requestData); !allowed {
				c.Response().Status = http.StatusTooManyRequests
				return echo.ErrTooManyRequests
			}

			// Go to next handler
			return next(c)
		}
	}
}

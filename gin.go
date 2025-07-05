package Daddy

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

func GinRateLimiter(redisClient *redis.Client, options ...services.Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a child context with timeout that inherits from the request context
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		// Check Redis availability first
		if err := redisClient.Ping(ctx).Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "Rate limiting service unavailable",
			})
			return
		}

		// Prepare request data
		requestData := &models.ReqData{
			UserIp:         c.ClientIP(),
			RequestAddress: c.Request.RequestURI,
		}

		// Create service instance
		service := services.NewService(db.CreateRedis(redisClient), options...)

		// Check rate limit
		if allowed := service.CheckAndStoreRate(ctx, requestData); !allowed {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// ez pz middleware done it's thing
		c.Next()
	}
}

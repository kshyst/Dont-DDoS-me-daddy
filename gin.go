package Daddy

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"github.com/redis/go-redis/v9"
	"time"
)

func GinRateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO : here we use empty context which is bad smell and should be using the passed context but the
		// passed context is timeouting soon
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		clientIP := c.ClientIP()
		requestedURL := c.Request.RequestURI

		requestData := &models.ReqData{
			UserIp:         clientIP,
			RequestAddress: requestedURL,
		}

		// check redis client availability
		if err := redisClient.Ping(c.Request.Context()).Err(); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Redis unavailable"})
			return
		}

		service := services.NewService(db.CreateRedis(redisClient))

		rateLimitRes := service.CheckAndStoreRate(c.Request.Context(), requestData)

		if rateLimitRes {
			c.Next()
		} else {
			c.Status(429)
			c.Abort()
		}

		fmt.Println("request processed in middleware")
	}
}

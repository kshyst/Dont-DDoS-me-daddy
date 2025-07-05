package Rate_Limiter

import (
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
		start := time.Now()
		defer func() {
			fmt.Printf("Middleware execution time: %v\n", time.Since(start))
		}()

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

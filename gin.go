package Rate_Limiter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kshyst/Dont-DDoS-me-daddy/db"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/models"
	"github.com/kshyst/Dont-DDoS-me-daddy/internal/services"
	"github.com/redis/go-redis/v9"
)

func GinRateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("processing request in middleware")

		clientIP := c.ClientIP()
		requestedURL := c.Request.RequestURI

		fmt.Println(clientIP, requestedURL)

		requestData := &models.ReqData{
			UserIp:         clientIP,
			RequestAddress: requestedURL,
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

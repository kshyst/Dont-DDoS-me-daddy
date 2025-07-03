package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/db"
	"main/internal/models"
	"main/internal/services"
)

func MyCustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("processing request in middleware")

		clientIP := c.ClientIP()
		requestedURL := c.Request.RequestURI

		requestData := &models.ReqData{
			UserIp:         clientIP,
			RequestAddress: requestedURL,
		}

		redis, _ := db.NewRedis()
		service := services.NewService(redis)

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

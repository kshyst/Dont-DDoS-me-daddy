package test

import (
	"context"
	"github.com/gin-gonic/gin"
	Daddy "github.com/kshyst/Dont-DDoS-me-daddy"
	"github.com/redis/go-redis/v9"
	"net/http"
	"testing"
	"time"
)

func TestGinMiddleware(t *testing.T) {
	const allowedRequestCount = 3

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	//test redis
	statusRedis := redisClient.Ping(context.Background())
	if statusRedis.Err() != nil {
		t.Errorf("Failed to start redis client %s", statusRedis.Err().Error())
	}

	r := gin.Default()
	r.Use(Daddy.GinRateLimiter(redisClient, Daddy.WithAllowedRequestCount(allowedRequestCount)))

	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
		return
	})

	go func() {
		err := r.Run(":8081")
		if err != nil {
			t.Errorf("Could not begin test because of gin router not starting: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	// making the url rate limiter activated
	for i := 0; i < allowedRequestCount; i++ {
		time.Sleep(2 * time.Second)
		_, err := http.Get("http://127.0.0.1:8081/test")
		if err != nil {
			t.Errorf("Could not execute GET request to check the middleware: %v", err)
		}
	}

	// send a request to be rate limited
	get, err := http.Get("http://127.0.0.1:8081/test")
	if err != nil {
		t.Errorf("Could not execute GET request to check the middleware: %v", err)
	}

	if get.StatusCode != http.StatusTooManyRequests {
		t.Errorf("The status code should be http.StatusTooManyRequests, instead got: %d", get.StatusCode)
	}
}

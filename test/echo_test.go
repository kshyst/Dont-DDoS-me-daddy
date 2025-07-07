package test

import (
	"context"
	Daddy "github.com/kshyst/Dont-DDoS-me-daddy"
	"github.com/labstack/echo"
	"github.com/redis/go-redis/v9"
	"net/http"
	"testing"
	"time"
)

func TestEchoMiddleware(t *testing.T) {
	const allowedRequestCount = 3

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	//test redis
	statusRedis := redisClient.Ping(context.Background())
	if statusRedis.Err() != nil {
		t.Errorf("Failed to start redis client %s", statusRedis.Err().Error())
	}

	e := echo.New()
	e.Use(Daddy.EchoMiddleware(redisClient, Daddy.WithAllowedRequestCount(allowedRequestCount)))

	e.GET("/test", func(context echo.Context) error {
		return context.String(200, "hello world")
	})

	go func() {
		e.Logger.Fatal(e.Start(":8082"))
	}()

	time.Sleep(3 * time.Second)

	// Send requests to trigger the ratelimiter
	for i := 0; i < allowedRequestCount; i++ {
		get, err := http.Get("http://127.0.0.1:8082/test")
		defer get.Body.Close()

		if err != nil {
			t.Errorf("Failed to send request: %s", err.Error())
			return
		} else if get.StatusCode != 200 {
			t.Errorf("Invalid status code. should be 200 but got: %d", get.StatusCode)
		} else {
			t.Logf("The status code is correct. Request number is %d", i)
		}
		time.Sleep(2 * time.Second)
	}

	// Send request that should be rate limited
	get, err := http.Get("http://127.0.0.1:8082/test")
	if err != nil {
		t.Errorf("Failed to send request: %s", err.Error())
		return
	} else if get.StatusCode != 429 {
		t.Errorf("Invalid status code. should be 429 but got: %d", get.StatusCode)
	} else {
		t.Logf("The status code is correct. Rate limiter working alright :)")
	}

}

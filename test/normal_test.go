package test

import (
	"context"
	Daddy "github.com/kshyst/Dont-DDoS-me-daddy"
	"github.com/redis/go-redis/v9"
	"net/http"
	"testing"
	"time"
)

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func TestRateLimiter(t *testing.T) {
	const allowedRequestCount = 5

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	//test redis
	statusRedis := redisClient.Ping(context.Background())
	if statusRedis.Err() != nil {
		t.Errorf("Failed to start redis client %s", statusRedis.Err().Error())
	}

	// Create and example handle func
	http.HandleFunc("/test", Daddy.RateLimiter(exampleHandler, redisClient, Daddy.WithAllowedRequestCount(allowedRequestCount)))

	// Start an http server
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			t.Errorf("Failed starting the test due to https server error: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	// Send requests to exceed the rate
	for i := 0; i < allowedRequestCount; i++ {
		get, err := http.Get("http://127.0.0.1:8080/test")
		if err != nil {
			t.Errorf("Failed to GET http://127.0.0.1:8080/test: %s", err.Error())
		} else if get.StatusCode != http.StatusOK {
			t.Errorf("The status code is wrong and should be 200 but got: %d", get.StatusCode)
		} else {
			t.Logf("The status code is correct. Request number is %d", i)
		}
		time.Sleep(2 * time.Second)
	}

	// Send another request and expect 429
	get, err := http.Get("http://127.0.0.1:8080/test")
	if err != nil {
		t.Errorf("Failed to GET http://127.0.0.1:8080/test For the final time: %s", err.Error())
	} else if get.StatusCode != http.StatusTooManyRequests {
		t.Errorf("The status code is wrong and should be 429 but got: %d", get.StatusCode)
	}
}

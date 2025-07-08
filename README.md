# Don't DDoS me daddy! 
![Go](https://img.shields.io/badge/Go-1.24-blue?logo=Go)
![Pr](https://img.shields.io/badge/PRs-welcome-lime?logo=Git)

Ddmd is an idiomatic rate limiting middleware which uses [Sliding Window Counter Algorithm](https://medium.com/@avocadi/rate-limiter-sliding-window-counter-7ec08dbe21d6) to control the request rate with various available options to customize.

It uses redis sorted list data type for storing request and holds them for the amount of expiration time which is customizable.

## Usage
Ddmd can be used as middleware in pure go handlerfuncs. I also made middlewares for some web frameworks for easier use. Beside middleware, it can be used as a microservice which you send request to and checks if it should be rate limited or not.

---

### HandlerFunc
You can use the `RateLimiter(next http.HandlerFunc, redisClient *redis.Client, options ...services.Option)` for normal HandlerFunc middleware usage like this:
```Go
http.HandleFunc(
  "/test", 
  Daddy.RateLimiter(exampleHandler, redisClient, Daddy.WithAllowedRequestCount(allowedRequestCount)),
)
```

In this example for wrapping up the `exampleHandler` with middleware we put it in the ratelimiter function and give some options if we want.

---

### Gin
[Gin](https://github.com/gin-gonic/gin) usage is pretty easy :)
```Go
r := gin.Default()
r.Use(Daddy.GinRateLimiter(redisClient, Daddy.WithAllowedRequestCount(allowedRequestCount)))
```
Where r is the gin engine (like `gin.Default()`).

---

### Echo
[Echo](https://github.com/labstack/echo) usage is also pretty easy :

```Go
e := echo.New()
e.Use(Daddy.EchoMiddleware(redisClient, Daddy.WithAllowedRequestCount(allowedRequestCount)))
```

> More detailed usages can be found in the test folder.

### Microservice
There is an option to run the main.go in cmd and also use the configs.yaml file to configure options.

You can then request to it with this format to check rate limitation:
```json
{
    "user_ip":"1.2.3.4",
    "request_address":"/url"
}
```

---

## Options

Options can be added to middlewares as parameters. It is shown in the examples.

### `WithWindowLength(seconds int)`
Length of which the sliding window should check for requests from sameuser.

### `WithAllowedRequestCount(count int)`
Count of the requests that can go through in the time window set with `WithWindowLength(seconds int)`.

### `WithExpiration(seconds int)`
Sets the expiration of the requests in redis (must be larger than window length).

### `WithRequestTimeout(seconds int)`
Timeout of the context that goes through middleware.
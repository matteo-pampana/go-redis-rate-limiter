package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/matteo-pampana/go-redis-rate-limiter/ratelimiter"
)

var (
	maxRequestsStr     = os.Getenv("RATE_LIMITER_MAX_REQUESTS")
	refreshIntervalStr = os.Getenv("RATE_LIMITER_REFRESH_INTERVAL")
	redisURI           = os.Getenv("REDIS_URI")
	serverPort         = os.Getenv("SERVER_PORT")
)

func main() {
	rateLimiter, err := initRateLimiter()
	if err != nil {
		panic(fmt.Errorf("error initializing rate limiter: %v", err))
	}

	r := gin.Default()

	r.GET("/rate-limiter", func(c *gin.Context) {
		name := c.Query("name")

		err := rateLimiter.CheckRequest(c.Request.Context(), []string{name})
		if err == ratelimiter.ErrTooManyRequests {
			c.JSON(429, gin.H{
				"message": err.Error(),
				"code":    429,
			})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
				"code":    500,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "hello " + name,
		})
	})

	r.Run(":" + serverPort)
}

func initRateLimiter() (*ratelimiter.RateLimiter, error) {
	maxRequests, err := strconv.ParseInt(maxRequestsStr, 10, 32)
	if err != nil {
		return nil, err
	}
	refreshInterval, err := time.ParseDuration(refreshIntervalStr)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})
	store := ratelimiter.NewStore(redisClient)
	config := ratelimiter.RateLimiterConfig{
		MaxRequests:     int(maxRequests),
		RefreshInterval: refreshInterval,
	}
	return ratelimiter.NewRateLimiter(store, config), nil
}

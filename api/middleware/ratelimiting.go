package middleware

import (
	"net/http"
	"strconv"
	"time"

	"URL-Shortner/constants"
	"URL-Shortner/log"
	"URL-Shortner/utils/cache"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		log.Sugar.Infof("Client IP: %s", clientIP)
		windowSize := strconv.ParseInt(viper.GetString("rateLimiting.windowSize")))
		maxRequests := viper.GetInt("rateLimiting.maxRequests")

		redisClient := cache.GetRedisCache().Client
		redisCount, err := redisClient.Incr(c, clientIP).Result()
		if err != nil {
			log.Sugar.Errorf("Redis GET error: %v", err)
			c.JSON(http.StatusInternalServerError, constants.ErrInternalServerError)
		}
		if redisCount == 1 {
			err := redisClient.Expire(c, clientIP, time.Duration(windowsize))
			if err != nil {
				log.Sugar.Errorf("Redis EXPIRE error: %v", err)
				c.JSON(http.StatusInternalServerError, constants.ErrInternalServerError)
			}

		}
		if redisCount > int64(maxRequests) {
			log.Sugar.Warnf("Rate limit exceeded for IP: %s", clientIP)
			c.JSON(http.StatusTooManyRequests, constants.ErrTooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}

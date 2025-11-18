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
		windowSize, err := strconv.ParseInt(viper.GetString("rateLimiting.windowSize"), 10, 64)
		if err != nil {
			log.Sugar.Errorf("Error parsing window size: %v", err)
			c.JSON(http.StatusInternalServerError, constants.ErrInternalServerError)
			return
		}
		maxRequests := viper.GetInt("rateLimiting.maxRequests")

		redisClient := cache.GetRedisCache().Client
		redisCount, err := redisClient.Incr(c, clientIP).Result()
		if err != nil {
			log.Sugar.Errorf("Redis GET error: %v", err)
			c.JSON(http.StatusInternalServerError, constants.ErrInternalServerError)
		}
		if redisCount == 1 {
			err := redisClient.Expire(c, clientIP, time.Duration(windowSize)*time.Millisecond).Err()
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

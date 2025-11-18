package cache

import (
	"context"
	"os"
	"time"

	"URL-Shortner/log"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type RedisCache struct {
	Client *redis.Client
}

var RedisClientInstance RedisCache

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
}

func InitRedisCache(ctx context.Context) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Sugar.Fatalf("Failed to parse Redis URL: %v", err)
	}
	RedisClientInstance = RedisCache{
		Client: redis.NewClient(opt),
	}
	err = RedisClientInstance.Client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
}

func GetRedisCache() *RedisCache {
	return &RedisClientInstance
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

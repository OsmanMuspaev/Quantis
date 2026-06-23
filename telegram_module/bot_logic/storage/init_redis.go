package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	AuthRedis *redis.Client
)

func InitRedis() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "tg_redis:6380"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	AuthRedis = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := AuthRedis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis connection failed: %v", err)
	}

	log.Println("Redis connected")
	return nil
}

func CloseRedis() {
	if AuthRedis != nil {
		AuthRedis.Close()
	}
	log.Println("Redis connection closed")
}

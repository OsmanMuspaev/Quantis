package storage

import (
	"github.com/redis/go-redis/v9"
)

// GetAuthRedis returns the Redis client instance.
func GetAuthRedis() *redis.Client {
	return AuthRedis
}

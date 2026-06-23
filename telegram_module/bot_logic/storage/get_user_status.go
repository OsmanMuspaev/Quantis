package storage

import (
	"context"
	"log"
	"strconv"
)

var ctx = context.Background()

func GetUserStatus(chatID int) string {
	id := strconv.Itoa(chatID)
	session, err := AuthRedis.Exists(ctx, id).Result()
	if err != nil {
		log.Printf("Redis error checking user %s: %v", id, err)
		return "unknown"
	}
	if session == 0 {
		return "unknown"
	}

	status, err := AuthRedis.HGet(ctx, id, "status").Result()
	if err != nil {
		log.Printf("Redis error getting status for %s: %v", id, err)
		return "unknown"
	}

	return status
}

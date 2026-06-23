package storage

import (
	"context"
	"strconv"
	"time"
)

func CreateUser(chatID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatIDStr := strconv.Itoa(chatID)
	entryToken := generateToken()

	err := AuthRedis.HSet(ctx, chatIDStr,
		"entry_token", entryToken,
		"status", "anonymous",
	).Err()
	if err != nil {
		return err
	}

	return AuthRedis.HSet(ctx, "ids", chatIDStr, chatIDStr).Err()
}

func UpdateUser(chatID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chatIDStr := strconv.Itoa(chatID)
	entryToken := generateToken()

	if err := AuthRedis.HSet(ctx, chatIDStr, "entry_token", entryToken).Err(); err != nil {
		return err
	}

	return AuthRedis.HSet(ctx, "ids", chatIDStr, chatIDStr).Err()
}

func GetEntryToken(chatID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return AuthRedis.HGet(ctx, strconv.Itoa(chatID), "entry_token").Result()
}

func GetRefreshToken(chatID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return AuthRedis.HGet(ctx, strconv.Itoa(chatID), "refresh_token").Result()
}

func generateToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

package update

import (
	"context"
	"log"

	"bot_logic/storage"
)

// CollectNotificationsFromAllUsers gathers notifications from main_module for all authorized users.
func CollectNotificationsFromAllUsers() (map[string][]Notification, error) {
	ctx := context.Background()

	userNotifications := make(map[string][]Notification)

	iter := storage.AuthRedis.Scan(ctx, 0, "*", 0).Iterator()

	for iter.Next(ctx) {
		chatID := iter.Val()

		userData, err := storage.AuthRedis.HGetAll(ctx, chatID).Result()
		if err != nil {
			log.Printf("Failed to HGETALL for chatID %s: %v", chatID, err)
			continue
		}

		status := userData["status"]
		token := userData["access_token"]

		if status != "authorized" || token == "" {
			continue
		}

		notifications, err := fetchUserNotifications(token)
		if err != nil {
			log.Printf("Failed to fetch notifications for chatID %s: %v", chatID, err)
			continue
		}

		if len(notifications) == 0 {
			continue
		}

		userNotifications[chatID] = notifications

		var ids []int
		for _, n := range notifications {
			ids = append(ids, n.ID)
		}
		confirmNotifications(token, ids)
	}
	if err := iter.Err(); err != nil {
		log.Printf("Redis SCAN iterator error: %v", err)
		return nil, err
	}
	return userNotifications, nil
}

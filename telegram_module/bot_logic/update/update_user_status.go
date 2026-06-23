package update

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"bot_logic/storage"
)

type TokenResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

func UpdateUserStatus() {
	ctx := context.Background()

	if err := storage.AuthRedis.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection error: %v", err)
		return
	}

	ids, err := storage.AuthRedis.HGetAll(ctx, "ids").Result()
	if err != nil {
		log.Printf("Error getting IDs from Redis: %v", err)
		return
	}

	if len(ids) == 0 {
		return
	}

	for id := range ids {
		entryToken, err := storage.AuthRedis.HGet(ctx, id, "entry_token").Result()
		if err != nil {
			log.Printf("Error getting entry_token for %s: %v", id, err)
			continue
		}

		url := "http://auth:8081/login/check?entry_token=" + entryToken

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Error checking status for %s: %v", id, err)
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error reading response for %s: %v", id, err)
			continue
		}

		var tokenResp TokenResponse
		if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
			log.Printf("Error parsing JSON for %s: %v", id, err)
			continue
		}

		switch tokenResp.Status {
		case "approved":
			err := storage.AuthRedis.HSet(ctx, id,
				"access_token", tokenResp.AccessToken,
				"refresh_token", tokenResp.RefreshToken,
				"user_id", tokenResp.UserID,
				"status", "authorized",
			).Err()
			if err != nil {
				log.Printf("Error setting authorized for %s: %v", id, err)
			}
			storage.AuthRedis.HDel(ctx, "ids", id)

		case "denied":
			storage.AuthRedis.HSet(ctx, id, "status", "unknown")
			storage.AuthRedis.HDel(ctx, "ids", id)

		case "pending":
			// Still waiting for user to enter code
		}
	}
}

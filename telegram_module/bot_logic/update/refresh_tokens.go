package update

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"bot_logic/storage"
)

type RefreshResponseStruct struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshUsersTokens() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			ctx := context.Background()

			iter := storage.AuthRedis.Scan(ctx, 0, "*", 0).Iterator()

			for iter.Next(ctx) {
				chatID := iter.Val()
				if chatID == "ids" {
					continue
				}

				id, _ := strconv.Atoi(chatID)
				userData, err := storage.AuthRedis.HGetAll(ctx, chatID).Result()
				if err != nil {
					log.Printf("Failed to HGETALL for chatID %s: %v", chatID, err)
					continue
				}

				status := userData["status"]
				if status != "authorized" {
					continue
				}

				refreshToken, err := storage.GetRefreshToken(id)
				if err != nil {
					log.Printf("Failed to get refresh token: %v", err)
					continue
				}

				data := map[string]string{"refresh_token": refreshToken}
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("Failed to create JSON: %v", err)
					continue
				}

				req, err := http.NewRequest("POST", "http://auth:8081/refresh", bytes.NewBuffer(jsonData))
				if err != nil {
					log.Printf("Failed to create request: %v", err)
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Request failed: %v", err)
					continue
				}

				bodyBytes, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					log.Printf("Failed to read response: %v", err)
					continue
				}

				var refreshResp RefreshResponseStruct
				if err := json.Unmarshal(bodyBytes, &refreshResp); err != nil {
					log.Printf("Failed to decode JSON: %v", err)
					continue
				}

				err = storage.AuthRedis.HSet(ctx, chatID,
					"access_token", refreshResp.AccessToken,
					"refresh_token", refreshResp.RefreshToken,
				).Err()
				if err != nil {
					log.Printf("Redis error: %v", err)
				}
			}
		}
	}()
}

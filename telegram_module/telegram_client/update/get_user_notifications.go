package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Notification struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type AllNotificationsResponse map[string][]Notification

func GetUserNotifications() {
	ticker := time.NewTicker(10 * time.Second)

	client := &http.Client{Timeout: 5 * time.Second}

	go func() {
		for range ticker.C {
			url := "http://tg_nginx/notifications"

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				continue
			}

			req.Header.Set("User-Agent", "NotificationService/1.0")
			req.Header.Set("Accept", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Failed to fetch notifications: %v", err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				log.Printf("Server returned non-OK status: %d", resp.StatusCode)
				continue
			}

			var allNotifications AllNotificationsResponse
			err = json.NewDecoder(resp.Body).Decode(&allNotifications)
			resp.Body.Close()
			if err != nil {
				log.Printf("Failed to decode response: %v", err)
				continue
			}

			if len(allNotifications) == 0 {
				continue
			}

			botToken := os.Getenv("BOT_TOKEN")
			if botToken == "" {
				log.Println("Warning: BOT_TOKEN not set")
				continue
			}

			for chatID, notifications := range allNotifications {
				for _, n := range notifications {
					sendTelegramNotification(botToken, chatID, n)
				}
			}
		}
	}()
}

func sendTelegramNotification(botToken, chatID string, notification Notification) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	fullText := fmt.Sprintf("<b>%s</b>\n\n%s", notification.Title, notification.Message)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       fullText,
		"parse_mode": "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error sending to Telegram API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned error %d", resp.StatusCode)
	}
}

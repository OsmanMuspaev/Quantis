package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Notification represents a single notification from the main module.
type Notification struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
}

type MainModuleResponse struct {
	Notifications []Notification `json:"notifications"`
}

func fetchUserNotifications(token string) ([]Notification, error) {
	req, err := http.NewRequest("GET", "http://main_module:18080/notification", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("main module returned status %d", resp.StatusCode)
	}

	var res MainModuleResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Notifications, nil
}

func confirmNotifications(token string, ids []int) {
	if len(ids) == 0 {
		return
	}

	body, _ := json.Marshal(map[string][]int{"ids": ids})
	req, _ := http.NewRequest("POST", "http://main_module:18080/notification/confirm-tg", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Confirm request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Confirm failed (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
}

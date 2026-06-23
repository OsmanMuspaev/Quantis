package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CodeResponse struct {
	Code string `json:"code"`
}

func CodeAuth(chatID int, entryToken string) (string, error) {
	body := map[string]string{"entry_token": entryToken}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "http://auth:8081/login?type=code", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("CodeAuth request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var codeResp CodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&codeResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return codeResp.Code, nil
}

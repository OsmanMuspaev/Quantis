package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"bot_logic/storage"
)

type VerifyCodeResponse struct {
	Status string `json:"status"`
}

func VerifyCode(chatID int, code string) error {
	refreshToken, err := storage.GetRefreshToken(chatID)
	if err != nil {
		return fmt.Errorf("refresh token not found: %v", err)
	}

	body := map[string]string{"code": code, "refresh_token": refreshToken}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "http://auth:8081/login/verify", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("VerifyCode request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var verifyResp VerifyCodeResponse
	if err := json.Unmarshal(bodyBytes, &verifyResp); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	if verifyResp.Status == "denied" {
		return fmt.Errorf("code verification denied")
	}
	return nil
}

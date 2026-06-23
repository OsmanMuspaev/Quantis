package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func VerifyCode(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, status string) {
	if status != "authorized" {
		m := tgbotapi.NewMessage(msg.Chat.ID, "You are not authorized. Please use /start first.")
		bot.Send(m)
		return
	}

	url := "http://tg_nginx/login/verify"

	data := map[string]string{"chat_id": fmt.Sprintf("%d", msg.Chat.ID), "code": msg.Text}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to create JSON: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return
	}

	m := tgbotapi.NewMessage(msg.Chat.ID, string(body))
	if _, err := bot.Send(m); err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	delete := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	bot.Send(delete)
}

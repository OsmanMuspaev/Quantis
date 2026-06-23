package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type UserInfoResponse struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
}

func GetUserInfoClientHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, status string) {
	if status != "authorized" {
		m := tgbotapi.NewMessage(msg.Chat.ID, "Please authorize first via /start")
		bot.Send(m)
		return
	}

	url := "http://tg_nginx/users/info?chat_id=" + strconv.FormatInt(msg.Chat.ID, 10)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	client := &http.Client{}
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

	var response UserInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return
	}

	text := "Username: " + response.FullName + " | User ID: " + response.UserID
	m := tgbotapi.NewMessage(msg.Chat.ID, text)
	if _, err := bot.Send(m); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

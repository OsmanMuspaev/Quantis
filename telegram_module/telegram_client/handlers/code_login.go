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

func CodeLoginBot(bot *tgbotapi.BotAPI, msg *tgbotapi.CallbackQuery, status string) {
	if status != "authorized" {
		url := "http://tg_nginx/login?type=code"

		data := map[string]int64{"chat_id": msg.Message.Chat.ID}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to create JSON: %v\n", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to create request: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request failed: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response: %v\n", err)
			return
		}

		var response struct {
			Code string `json:"code"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Failed to parse JSON: %v\n", err)
			return
		}

		messageText := fmt.Sprintf(
			"*Code for login:*\n\n"+
				"```\n%s\n```\n\n"+
				"*How to use:*\n"+
				"1. Open the app on an authorized device\n"+
				"2. Use /verifycode command\n"+
				"3. Enter the code\n\n"+
				"*Valid for:* 60 seconds\n"+
				"*Do not share this code!*",
			response.Code,
		)

		m := tgbotapi.NewMessage(msg.Message.Chat.ID, messageText)
		m.ParseMode = "Markdown"
		bot.Send(m)
	} else {
		m := tgbotapi.NewMessage(msg.Message.Chat.ID, "You are already authorized!")
		bot.Send(m)
	}
	delete := tgbotapi.NewDeleteMessage(msg.Message.Chat.ID, msg.Message.MessageID)
	bot.Send(delete)
}

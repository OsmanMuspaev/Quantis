package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"telegram_client/state"
)

type GetUserStatusResponse struct {
	Status string `json:"status"`
}

func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, update *tgbotapi.Update) {
	url := "http://tg_nginx/get_user_status"

	body := map[string]int64{"chat_id": msg.Chat.ID}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	stat, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}
	defer stat.Body.Close()

	if stat.StatusCode != http.StatusOK {
		log.Printf("Server returned error: %s", stat.Status)
		return
	}

	var status GetUserStatusResponse
	if err := json.NewDecoder(stat.Body).Decode(&status); err != nil {
		log.Printf("Failed to decode status: %v", err)
		return
	}

	switch msg.Text {
	case "/start":
		StartBot(bot, msg)
	case "/verifycode":
		if status.Status == "authorized" {
			m := tgbotapi.NewMessage(msg.Chat.ID, "Enter your verification code:")
			bot.Send(m)
			state.SetUserState(msg.Chat.ID, "awaiting_code")
		} else {
			m := tgbotapi.NewMessage(msg.Chat.ID, "Please authorize first via /start")
			bot.Send(m)
		}
	case "/getuserinfo":
		if status.Status == "authorized" {
			GetUserInfoClientHandler(bot, msg, status.Status)
		} else {
			m := tgbotapi.NewMessage(msg.Chat.ID, "Please authorize first via /start")
			bot.Send(m)
		}
	case "/getuserprofile":
		if status.Status == "authorized" {
			GetUserProfileClientHandler(bot, msg, status.Status)
		} else {
			m := tgbotapi.NewMessage(msg.Chat.ID, "Please authorize first via /start")
			bot.Send(m)
		}
	default:
		if state.GetUserState(msg.Chat.ID) == "awaiting_code" {
			VerifyCode(bot, msg, status.Status)
			state.SetUserState(msg.Chat.ID, "idle")
		} else {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Unknown command"))
		}
	}
}

func HandleCallback(bot *tgbotapi.BotAPI, msg *tgbotapi.CallbackQuery, update *tgbotapi.Update) {
	url := "http://tg_nginx/get_user_status"

	body := map[string]int64{"chat_id": msg.Message.Chat.ID}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	stat, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return
	}
	defer stat.Body.Close()

	if stat.StatusCode != http.StatusOK {
		log.Printf("Server returned error: %s", stat.Status)
		return
	}

	var status GetUserStatusResponse
	if err := json.NewDecoder(stat.Body).Decode(&status); err != nil {
		log.Printf("Failed to decode status: %v", err)
		return
	}

	callback := update.CallbackQuery
	data := callback.Data

	answer := tgbotapi.NewCallback(callback.ID, "")
	bot.AnswerCallbackQuery(answer)

	if data == "login" {
		LoginBot(bot, msg, status.Status)
	} else if data == "github_login" {
		GithubLoginBot(bot, msg, status.Status)
	} else if data == "yandex_login" {
		YandexLoginBot(bot, msg, status.Status)
	} else if data == "code_au" {
		CodeLoginBot(bot, msg, status.Status)
	}
}

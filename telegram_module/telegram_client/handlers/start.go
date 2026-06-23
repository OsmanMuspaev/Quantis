package handlers

import (
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func StartBot(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	reply := "Welcome! To continue, please authorize:"

	loginBtn := tgbotapi.NewInlineKeyboardButtonData("Login / Authorize", "login")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(loginBtn),
	)

	mesg := tgbotapi.NewMessage(msg.Chat.ID, reply)
	mesg.ReplyMarkup = keyboard

	if _, err := bot.Send(mesg); err != nil {
		log.Printf("Failed to send start message: %v", err)
	}
}

package handlers

import (
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func LoginBot(bot *tgbotapi.BotAPI, msg *tgbotapi.CallbackQuery, status string) {
	m := tgbotapi.NewMessage(msg.Message.Chat.ID, "Choose a login method:")

	githubLogin := tgbotapi.NewInlineKeyboardButtonData("Github OAuth", "github_login")
	yandexLogin := tgbotapi.NewInlineKeyboardButtonData("Yandex OAuth", "yandex_login")
	codeAuth := tgbotapi.NewInlineKeyboardButtonData("Code Auth", "code_au")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(githubLogin),
		tgbotapi.NewInlineKeyboardRow(yandexLogin),
		tgbotapi.NewInlineKeyboardRow(codeAuth),
	)
	m.ReplyMarkup = keyboard

	if _, err := bot.Send(m); err != nil {
		log.Printf("Failed to send login options: %v", err)
	}
	delete := tgbotapi.NewDeleteMessage(msg.Message.Chat.ID, msg.Message.MessageID)
	bot.Send(delete)
}

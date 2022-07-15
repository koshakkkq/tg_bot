package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var users_keybords map[string]tgbotapi.ReplyKeyboardMarkup

func load_keybords() {
	users_keybords = make(map[string]tgbotapi.ReplyKeyboardMarkup)
	users_keybords["default"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить ссылку"),
			tgbotapi.NewKeyboardButton("Получить ссылку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Мои кредиты")),
	)
	users_keybords["back"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
}

func return_to_default(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
	msg.ReplyMarkup = users_keybords["default"]

	if _, err := bot.Send(msg); err != nil {
		errorLog.Println(err)
		return
	}
	id := update.Message.From.ID
	if len(users_status[id]) > 1 {
		users_status[id] = users_status[id][:1]
	}
	users_status[id][0] = 1
}

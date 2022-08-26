package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (main *Admin) Load_Admin_keybords() {
	main.Admin_keybords = make(map[string]tgbotapi.ReplyKeyboardMarkup)
	main.Admins_status = make(map[int64][]int64)
	main.Cur_room = make(map[int64]string)
	main.Admin_keybords["default"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить ссылку"),
			tgbotapi.NewKeyboardButton("Получить ссылку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Админ"),
		),
	)
	main.Admin_keybords["back"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["panel"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Статистика по пользователям"),
			tgbotapi.NewKeyboardButton("Статистика по админам"),
			tgbotapi.NewKeyboardButton("Статистика по нику"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Увеличить баланс"),
			tgbotapi.NewKeyboardButton("Whitelist"),
			tgbotapi.NewKeyboardButton("Комнаты"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["forward_prev"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Следующая страница"),
			tgbotapi.NewKeyboardButton("Предыдущая страница"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["forward"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Следующая страница"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["prev"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Предыдущая страница"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["whitelist"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Ожидающие"),
			tgbotapi.NewKeyboardButton("Добавить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["rooms"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Управление"),
			tgbotapi.NewKeyboardButton("Все комнаты"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	main.Admin_keybords["room_page"] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавть пользователей"),
			tgbotapi.NewKeyboardButton("Рассылка"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Получить ники пользователей"),
			tgbotapi.NewKeyboardButton("Получить ссылки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
}
func (main *Admin) return_to_default_admin(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
	msg.ReplyMarkup = main.Admin_keybords["default"]

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	id := update.Message.From.ID
	main.Admins_status[id] = []int64{}
}

func (main *Admin) return_to_panel(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ, Выберите действие")
	msg.ReplyMarkup = main.Admin_keybords["panel"]

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	id := update.Message.From.ID
	main.Admins_status[id] = main.Admins_status[id][:2]
	main.Admins_status[id][1] = 2
}

func (main *Admin) return_to_whitelist_panel(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ, Выберите действие(вайтлист)")
	msg.ReplyMarkup = main.Admin_keybords["whitelist"]

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	id := update.Message.From.ID
	main.Admins_status[id] = main.Admins_status[id][:4]
	main.Admins_status[id][3] = 1
}

func (main *Admin) return_to_rooms_panel(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ, Выберите действие(комнаты)")
	msg.ReplyMarkup = main.Admin_keybords["rooms"]

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	id := update.Message.From.ID
	main.Admins_status[id] = main.Admins_status[id][:4]
	main.Admins_status[id][3] = 1
}

func (main *Admin) return_to_rooms_page(update *tgbotapi.Update) {
	id := update.Message.From.ID
	room_name := main.Cur_room[id]
	main.Cur_room[id] = room_name
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команата: "+room_name)
	msg.ReplyMarkup = main.Admin_keybords["room_page"]

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	main.Admins_status[id][3] = 3
}

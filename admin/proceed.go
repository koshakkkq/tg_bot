package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (main *Admin) Admin_proceed(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) == 0 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	if update.Message.Text == "Получить ссылку" || update.Message.Text == "Добавить ссылку" || update.Message.Text == "Админ" {
		main.Admins_status[id][0] = 1
	}
	switch main.Admins_status[id][0] {
	case 0: // Получить главную клавиатуру
		main.main_page_admin(update)
		return
	case 1: // Выбор действия
		main.choice_activity_admin(update)
		return
	}
}
func (main *Admin) choice_activity_admin(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 2 {
		main.Admins_status[id] = append(main.Admins_status[id], -1)
	}
	switch main.Admins_status[id][1] {
	case -1: //Ничего не выбрано
		switch update.Message.Text {
		case "Добавить ссылку":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ссылку")
			msg.ReplyMarkup = main.Admin_keybords["back"]

			if _, err := main.Bot.Send(msg); err != nil {
				main.return_to_default_admin(update)
				main.ErrorLog.Println(err)
				return
			}
			main.Admins_status[id][1] = 1
			return
		case "Получить ссылку":
			main.get_ref(update)
			return

		case "Админ":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ, Выберите действие")
			msg.ReplyMarkup = main.Admin_keybords["panel"]
			main.Admins_status[id][1] = 2
			if _, err := main.Bot.Send(msg); err != nil {
				main.return_to_default_admin(update)
				main.ErrorLog.Println(err)
				return
			}
			return
		default:
			main.return_to_default_admin(update)
			return
		}
	case 1: // Выбрано добавление ссылки
		switch update.Message.Text {
		case "Назад":
			main.return_to_default_admin(update)
			return
		default:
			main.add_ref(update)
			return
		}
	case 2: // Выбрана Админская панель
		main.proceed_admin_panel(update)
		return
	}
}
func (main *Admin) proceed_admin_panel(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 3 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	switch main.Admins_status[id][2] {
	case 0: // Ничего не выбрано
		switch update.Message.Text {
		case "Статистика по пользователям":
			main.proceed_stat(update, false)
			main.Admins_status[id][2] = 1
			return
		case "Статистика по админам":
			main.proceed_stat(update, true)
			main.Admins_status[id][2] = 2
			return
		case "Статистика по нику":
			main.proceed_name_stat(update)
			main.Admins_status[id][2] = 3
			return
		case "Увеличить баланс":
			main.proceed_balance_update(update)
			main.Admins_status[id][2] = 4
			return
		case "Whitelist":
			main.proceed_whitelist(update)
			main.Admins_status[id][2] = 5
			return
		case "Комнаты":
			main.proceed_rooms(update)
			main.Admins_status[id][2] = 6
		case "Назад":
			main.return_to_default_admin(update)
			return
		default:
			main.return_to_panel(update)
			return
		}
	case 1: // Стата по пользователям
		main.proceed_stat(update, false)
		return
	case 2: // стата по админам
		main.proceed_stat(update, true)
		return
	case 3: // стата по нику
		main.proceed_name_stat(update)
		return
	case 4: // Добавить баланс
		main.proceed_balance_update(update)
		return
	case 5: // whitelist
		main.proceed_whitelist(update)
		return
	case 6: // комнаты
		main.proceed_rooms(update)
		return
	}

}
func (main *Admin) proceed_name_stat(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 4 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	switch main.Admins_status[id][3] {
	case 0: //Ждём ник
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введитие ник")
		msg.ReplyMarkup = main.Admin_keybords["back"]
		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		main.Admins_status[id][3] = 1
	case 1:
		switch update.Message.Text {
		case "Назад":
			main.return_to_panel(update)
			return
		default:
			main.get_infrom_by_name(update)
		}

	}
}
func (main *Admin) proceed_stat(update *tgbotapi.Update, admin bool) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 4 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	number_of_users, err := main.get_users_count(admin)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	switch update.Message.Text {
	case "Следующая страница":
		if main.Admins_status[id][3]*10 < number_of_users {
			main.Admins_status[id][3]++
		}
	case "Предыдущая страница":
		if main.Admins_status[id][3] > 0 {
			main.Admins_status[id][3]--
		}
	case "Назад":
		main.return_to_panel(update)
		return
	}
	main.get_stat(update, admin)
}
func (main *Admin) proceed_balance_update(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 4 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	switch main.Admins_status[id][3] {
	case 0: //Ждём ник
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введитие ник и значение на которое изменится баланс(через пробел)")
		msg.ReplyMarkup = main.Admin_keybords["back"]
		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		main.Admins_status[id][3] = 1
	case 1:
		switch update.Message.Text {
		case "Назад":
			main.return_to_panel(update)
			return
		default:
			main.update_balance(update)
			return
		}

	}

}
func (main *Admin) proceed_whitelist(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 4 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	switch main.Admins_status[id][3] {
	case 0: //Ждём что-то
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ, выберите действие")
		msg.ReplyMarkup = main.Admin_keybords["whitelist"]
		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		main.Admins_status[id][3] = 1
		return
	case 1:
		switch update.Message.Text {
		case "Ожидающие":
			main.get_waiting(update)
			return
		case "Добавить":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ники через пробел")
			msg.ReplyMarkup = main.Admin_keybords["back"]
			if _, err := main.Bot.Send(msg); err != nil {
				main.Admins_status[id] = []int64{}
				main.ErrorLog.Println(err)
				return
			}
			main.Admins_status[id][3] = 2
			return
		case "Назад":
			main.return_to_panel(update)
			return
		}
	case 2:
		switch update.Message.Text {
		case "Назад":
			main.return_to_whitelist_panel(update)
			return
		default:
			main.add_to_whitelist(update)
			return
		}

	}
}

func (main *Admin) proceed_rooms(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(main.Admins_status[id]) < 4 {
		main.Admins_status[id] = append(main.Admins_status[id], 0)
	}
	switch main.Admins_status[id][3] {
	case 0: // Ждём ввода
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ, выберите действие")
		msg.ReplyMarkup = main.Admin_keybords["rooms"]
		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		main.Admins_status[id][3] = 1
		return
	case 1: // Выбор
		switch update.Message.Text {
		case "Назад":
			main.return_to_panel(update)
			return
		case "Управление":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите имя комнаты")
			msg.ReplyMarkup = main.Admin_keybords["back"]
			if _, err := main.Bot.Send(msg); err != nil {
				main.Admins_status[id] = []int64{}
				main.ErrorLog.Println(err)
				return
			}
			main.Admins_status[id][3] = 2
			return
		case "Все комнаты":
			main.get_all_rooms(update)
		}
	case 2: // Ожидание Имени комнаты
		switch update.Message.Text {
		case "Назад":
			main.return_to_rooms_panel(update)
			return
		default:
			main.set_cur_room(update)
			main.Admins_status[id][3] = 3
			return
		}
	case 3:
		switch update.Message.Text {
		case "Получить ники пользователей":
			main.get_room_users_full_info(update)
			return
		case "Добавть пользователей":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ники через пробел")
			msg.ReplyMarkup = main.Admin_keybords["back"]
			if _, err := main.Bot.Send(msg); err != nil {
				main.Admins_status[id] = []int64{}
				main.ErrorLog.Println(err)
				return
			}
			main.Admins_status[id][3] = 4
			return
		case "Рассылка":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите текст расслыки")
			msg.ReplyMarkup = main.Admin_keybords["back"]
			if _, err := main.Bot.Send(msg); err != nil {
				main.Admins_status[id] = []int64{}
				main.ErrorLog.Println(err)
				return
			}
			main.Admins_status[id][3] = 5
			return
		case "Получить ссылки":
			main.get_ref_by_room(update)
			return
		case "Назад":
			main.return_to_rooms_panel(update)
			return
		}
	case 4: // Ждём ники
		switch update.Message.Text {
		case "Назад":
			main.return_to_rooms_page(update)
			return
		default:
			main.add_users(update)
			return
		}
	case 5:
		switch update.Message.Text { //Ждём текст расслыки
		case "Назад":
			main.return_to_rooms_page(update)
			return
		default:
			main.create_mailing(update)
			return
		}
	}
}

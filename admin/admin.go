package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"strconv"
	"strings"
)

func (main *Admin) get_infrom_by_name(update *tgbotapi.Update) {
	name := update.Message.Text
	id := update.Message.From.ID
	user, err := main.get_user_by_name(name)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	if user == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пользователь с ником "+name+", не найден!")

		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		return
	}
	was := make(map[string]int)
	for _, el := range user.Refs {
		was[el] = 1
	}
	real_ref_was := make([]string, 0)
	for _, el := range user.Refs_was {
		if was[el] == 0 {
			real_ref_was = append(real_ref_was, el)
		}
	}
	user.Refs_was = real_ref_was
	tmpl, err := template.New("data").Parse("<a href='tg://user?id={{.Id}}'>{{.Name}}</a>\n" +
		"Баланс: {{.Balance}}\n" +
		"Админ:{{.Admin}}")
	msg_str := new(strings.Builder)

	err = tmpl.Execute(msg_str, user)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_str.String())

	msg.ParseMode = "HTML"
	if _, err := main.Bot.Send(msg); err != nil {
		main.Admins_status[id] = []int64{}
		main.ErrorLog.Println(err)
		return
	}
	refs_msgs := []string{"Добавленные ссылки:\n"}
	j := 0
	for _, el := range user.Refs {
		if len(el)+len(refs_msgs[j]) > 4000 {
			refs_msgs = append(refs_msgs, "")
			j++
		}
		refs_msgs[j] += el + "\n"
	}

	refs_was_msgs := []string{"Открытые ссылки:\n"}
	j = 0
	for _, el := range user.Refs_was {
		if len(el)+len(refs_was_msgs[j]) > 4000 {
			refs_was_msgs = append(refs_was_msgs, "")
			j++
		}
		refs_was_msgs[j] += el + "\n"
	}
	for _, el := range refs_msgs {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, el)

		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
	}

	for _, el := range refs_was_msgs {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, el)

		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
	}
	main.return_to_panel(update)

}

type user_tmpl struct {
	Name      string
	Id        int64
	Refs_open int
	Refs_aded int
	Balance   int64
	Admin     bool
}

func (main *Admin) get_stat(update *tgbotapi.Update, admin bool) {
	id := update.Message.From.ID
	option := options.Find()
	option.SetSort(bson.M{"name": 1})
	option.SetSkip(main.Admins_status[id][3] * 10)
	option.SetLimit(10)

	cur, err := main.Db_users.Find(context.TODO(), bson.M{"admin": admin}, option)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	var msg_text string
	for cur.Next(context.TODO()) {
		var user user_struct
		err = cur.Decode(&user)
		if err != nil {
			main.proceed_err(update, err)
			return
		}
		data := user_tmpl{
			Name:      user.Name,
			Id:        user.Id,
			Refs_open: len(user.Refs_was) - len(user.Refs),
			Refs_aded: len(user.Refs),
		}
		tmpl, err := template.New("data").Parse("<a href='tg://user?id={{.Id}}'>{{.Name}}</a>\n" +
			"            Добавленно ссылок: {{.Refs_aded}}\n" +
			"            Запрошено ссылко: {{.Refs_open}}\n" +
			"            Баланс: {{.Balance}}\n")
		if err != nil {
			main.proceed_err(update, err)
			return
		}
		msg_str := new(strings.Builder)

		err = tmpl.Execute(msg_str, data)
		if err != nil {
			main.proceed_err(update, err)
			return
		}
		msg_text += msg_str.String()
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_text)
	msg.ParseMode = "HTML"

	buttons := []tgbotapi.KeyboardButton{}
	number_of_users, err := main.get_users_count(admin)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	if main.Admins_status[id][3] != 0 {
		buttons = append(buttons, tgbotapi.NewKeyboardButton("Предыдущая страница"))
	}
	if (main.Admins_status[id][3]+1)*10 < number_of_users {
		buttons = append(buttons, tgbotapi.NewKeyboardButton("Следующая страница"))
	}

	keyboard := tgbotapi.NewReplyKeyboard(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	))
	msg.ReplyMarkup = keyboard
	if _, err := main.Bot.Send(msg); err != nil {
		main.Admins_status[id] = []int64{}
		main.ErrorLog.Println(err)
		return
	}
	return
}
func (main *Admin) update_balance(update *tgbotapi.Update) {
	id := update.Message.From.ID
	words := strings.Fields(update.Message.Text)
	if len(words) != 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Проверь корректность запроса, он должен иметь вид(без скобок)\n [ник] [значение на которое изменяется баланс]")
		if _, err := main.Bot.Send(msg); err != nil {
			main.ErrorLog.Println(err)
			return
		}
	}

	user, err := main.get_user_by_name(words[0])
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	if user == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пользователь с ником "+words[0]+", не найден!")

		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		return
	}
	dif, err := strconv.ParseInt(words[1], 10, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Проверь корректность запроса, он должен иметь вид(без скобок)\n [ник] [значение на которое изменяется баланс]")

		if _, err := main.Bot.Send(msg); err != nil {
			main.Admins_status[id] = []int64{}
			main.ErrorLog.Println(err)
			return
		}
		return
	}

	_, err = main.Db_users.UpdateOne(context.TODO(), bson.M{"id": user.Id}, bson.M{"$inc": bson.M{"balance": dif}})
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	tmpl, err := template.New("data").Parse("Баланс пользователя <a href='tg://user?id={{.Id}}'>{{.Name}}</a> успешно изменён")
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msg_str := new(strings.Builder)

	err = tmpl.Execute(msg_str, user)
	if err != nil {
		main.proceed_err(update, err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_str.String())
	msg.ParseMode = "HTML"
	if _, err := main.Bot.Send(msg); err != nil {
		main.Admins_status[id] = []int64{}
		main.ErrorLog.Println(err)
		return
	}
	main.return_to_panel(update)
}

func (main *Admin) get_users_count(admin bool) (int64, error) {
	cnt, err := main.Db_users.CountDocuments(context.TODO(), bson.M{"admin": admin})
	return cnt, err
}
func (main *Admin) main_page_admin(update *tgbotapi.Update) {
	main.return_to_default_admin(update)
}
func (main *Admin) proceed_err(update *tgbotapi.Update, err error) {
	main.ErrorLog.Println(err)
	id := update.Message.From.ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка на сервере")
	msg.ReplyMarkup = main.Admin_keybords["default"]
	main.Admins_status[id] = []int64{}
	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	return
}

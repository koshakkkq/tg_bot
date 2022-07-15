package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"strings"
)

func (main *Admin) set_cur_room(update *tgbotapi.Update) {
	id := update.Message.From.ID
	room_name := update.Message.Text
	main.Cur_room[id] = room_name
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Перехожу к комнате "+room_name)
	msg.ReplyMarkup = main.Admin_keybords["room_page"]

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
}
func (main *Admin) get_room_users_full_info(update *tgbotapi.Update) {
	id := update.Message.From.ID
	room_name := main.Cur_room[id]
	cur, err := main.Db_users.Find(context.TODO(), bson.M{"room": room_name})
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	j := 0
	msgs := []string{""}
	for cur.Next(context.TODO()) {
		var user user_struct
		cur.Decode(&user)
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
		if len(msg_str.String())+len(msgs) > 4000 {
			msgs = append(msgs, "")
			j++
		}
		msgs[j] += msg_str.String()
	}
	for _, el := range msgs {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, el)
		msg.ParseMode = "HTML"

		if _, err := main.Bot.Send(msg); err != nil {
			main.ErrorLog.Println(err)
			return
		}
	}
}

func (main *Admin) add_users(update *tgbotapi.Update) {
	id := update.Message.From.ID
	room_name := main.Cur_room[id]
	words := strings.Fields(update.Message.Text)
	real_words := []string{}
	for _, el := range words {
		if el[0] == '@' {
			real_words = append(real_words, el)
		}
	}
	_, err := main.Db_users.UpdateMany(context.TODO(), bson.M{"name": bson.M{"$in": real_words}}, bson.M{"$set": bson.M{"room": room_name}})
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно")

	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	main.return_to_rooms_page(update)
}
func (main *Admin) get_ref_by_room(update *tgbotapi.Update) {
	id := update.Message.From.ID
	room_name := main.Cur_room[id]
	cur, err := main.Db_refs.Find(context.TODO(), bson.M{"room": room_name})
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msgs := []string{""}
	j := 0
	for cur.Next(context.TODO()) {
		var ref ref_struct
		err := cur.Decode(&ref)
		if err != nil {
			main.proceed_err(update, err)
			return
		}
		if len(msgs)+len(ref.Ref) > 4000 {
			msgs = append(msgs, "")
		}
		msgs[j] += ref.Ref + "\n"
	}
	for _, el := range msgs {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, el)

		if _, err := main.Bot.Send(msg); err != nil {
			main.ErrorLog.Println(err)
			return
		}
	}
}

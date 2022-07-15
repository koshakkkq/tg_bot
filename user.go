package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"strconv"
	"strings"
)

func user_proceed(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(users_status[id]) == 0 {
		users_status[id] = append(users_status[id], 0)
	}
	if update.Message.Text == "Получить ссылку" || update.Message.Text == "Добавить ссылку" {
		users_status[id][0] = 1
	}
	switch users_status[id][0] {
	case 0: // Получить главную клавиатуру
		main_page(update)
		return
	case 1: // Выбор действия
		choice_activity(update)
		return
	}
}

func main_page(update *tgbotapi.Update) {
	return_to_default(update)
}

func choice_activity(update *tgbotapi.Update) {
	id := update.Message.From.ID
	if len(users_status[id]) < 2 {
		users_status[id] = append(users_status[id], -1)
	}
	switch users_status[id][1] {
	case -1:
		switch update.Message.Text {
		case "Добавить ссылку":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ссылку")
			msg.ReplyMarkup = users_keybords["back"]

			if _, err := bot.Send(msg); err != nil {
				users_status[id] = []int64{}
				errorLog.Println(err)
				return
			}
			users_status[id][1] = 1
			return
		case "Получить ссылку":
			get_ref(update)
			return

		case "Мои кредиты":
			get_my_credits(update)
		}
	case 1: // добавление ссылки
		switch update.Message.Text {
		case "Назад":
			return_to_default(update)
		default:
			add_ref(update)
		}

	}

}

func add_ref(update *tgbotapi.Update) {
	id := update.Message.From.ID

	balance, err := get_balance(id)
	if err != nil {
		proceed_err(update, err)
		return
	}

	if balance < 10 {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас недостаточно кредитов нужно 10, у вас: "+strconv.FormatInt(balance, 10))

		if _, err := bot.Send(msg); err != nil {
			users_status[id] = []int64{}
			errorLog.Println(err)
			return
		}

		return_to_default(update)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
		msg.ReplyMarkup = users_keybords["default"]

		if _, err := bot.Send(msg); err != nil {
			users_status[id] = []int64{}
			errorLog.Println(err)
			return
		}
		return
	}
	ref := update.Message.Text

	update_db := bson.M{
		"$inc": bson.M{
			"balance": -10,
		},
	}
	_, err = db_users.UpdateOne(context.TODO(), bson.M{"id": id}, update_db)
	if err != nil {
		proceed_err(update, err)
		return
	}

	update_db = bson.M{
		"$push": bson.M{
			"refs":     ref,
			"refs_was": ref,
		},
	}

	_, err = db_users.UpdateOne(context.TODO(), bson.M{"id": id}, update_db)
	if err != nil {
		proceed_err(update, err)
		return
	}
	room_name, err := get_room(id)
	if err != nil {
		proceed_err(update, err)
		return
	}

	ref_db_update := ref_struct{
		Ref:  ref,
		From: int64(id),
		Room: room_name,
	}
	_, err = db_refs.InsertOne(context.TODO(), ref_db_update)
	if err != nil {
		proceed_err(update, err)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно добавлена ссылка")

	if _, err := bot.Send(msg); err != nil {
		users_status[id] = []int64{}
		errorLog.Println(err)
		return
	}
	return_to_default(update)

}

type ref_tmpl struct {
	Ref     string
	Balance int64
}

func get_ref(update *tgbotapi.Update) {
	id := update.Message.From.ID

	var user user_struct

	err := db_users.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		proceed_err(update, err)
		return
	}
	filter := bson.M{"ref": bson.M{"$nin": user.Refs_was}}
	var ref ref_struct
	err = db_refs.FindOne(context.TODO(), filter).Decode(&ref)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нет ссылки, которую вам возможно выдать, попробуй позже")
			if _, err := bot.Send(msg); err != nil {
				users_status[id] = []int64{}
				errorLog.Println(err)
				return
			}
			return
		} else {
			proceed_err(update, err)
			return
		}
	}
	_, err = db_users.UpdateOne(context.TODO(), bson.M{"id": ref.From}, bson.M{"$inc": bson.M{"balance": 1}})

	if err != nil {
		proceed_err(update, err)
		return
	}

	update_db := bson.M{
		"$inc": bson.M{
			"balance": 1,
		},
		"$push": bson.M{
			"refs_was": ref.Ref,
		},
	}
	_, err = db_users.UpdateOne(context.TODO(), bson.M{"id": id}, update_db)
	if err != nil {
		proceed_err(update, err)
		return
	}
	tmpl_data := ref_tmpl{
		Ref:     ref.Ref,
		Balance: user.Balance + 1,
	}

	tmpl, err := template.New("data").Parse("<a href='{{.Ref}}'>Ссылка</a>\n<a>Твой баланс: {{.Balance}}</a>")
	if err != nil {
		proceed_err(update, err)
		return
	}
	msg_str := new(strings.Builder)

	err = tmpl.Execute(msg_str, tmpl_data)
	if err != nil {
		proceed_err(update, err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_str.String())
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		users_status[id] = []int64{}
		errorLog.Println(err)
		return
	}
	return
}

func get_my_credits(update *tgbotapi.Update) {
	balance, err := get_balance(update.Message.From.ID)
	if err != nil {
		proceed_err(update, err)
		return
	}
	user_balance := strconv.FormatInt(balance, 10)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Твой баланс: "+user_balance)
	if _, err := bot.Send(msg); err != nil {
		users_status[update.Message.From.ID] = []int64{}
		errorLog.Println(err)
		return
	}
}

func get_balance(id int64) (int64, error) {
	var user user_struct
	err := db_users.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func get_room(id int64) (string, error) {
	var user user_struct
	err := db_users.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.Room, nil
}

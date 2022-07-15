package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func procced_super_admin(update *tgbotapi.Update) {
	words := strings.Fields(update.Message.Text)
	if len(words) != 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Никнейм не введён")

		if _, err := bot.Send(msg); err != nil {
			errorLog.Println(err)
			return
		}
		return
	}

	name := words[1]

	if name[0] != '@' {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Никнейм должен начинаться с символа @")

		if _, err := bot.Send(msg); err != nil {
			errorLog.Println(err)
			return
		}
		return
	}
	if words[0] == "!add" {
		addAdmin(update)
	} else if words[0] == "!del" {
		delete_admin(update)
	}
}
func addAdmin(update *tgbotapi.Update) {
	words := strings.Fields(update.Message.Text)

	name := words[1]

	user, err := get_user_by_name(name)
	if err != nil {
		proceed_err(update, err)
		return
	}
	if user == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пользователь с ником "+name+", не найден!\nПользователь "+name+" должен написть боту!")

		if _, err := bot.Send(msg); err != nil {
			errorLog.Println(err)
			return
		}
		return
	}

	err = add_admin_privilege(user.Id)
	if err != nil {
		proceed_err(update, err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно добавлен админ "+name)

	users_status[user.Id] = []int64{}
	admins_status[user.Id] = []int64{}
	if _, err := bot.Send(msg); err != nil {
		errorLog.Println(err)
		return
	}
}
func delete_admin(update *tgbotapi.Update) {
	update_admin_db()
	words := strings.Fields(update.Message.Text)

	name := words[1]

	user, err := get_user_by_name(name)

	if err != nil {
		proceed_err(update, err)
		return
	}
	err = delete_admin_by_Id(user.Id)
	if err != nil {
		proceed_err(update, err)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно удалён админ "+name)

	users_status[user.Id] = []int64{}
	admins_status[user.Id] = []int64{}
	if _, err := bot.Send(msg); err != nil {
		errorLog.Println(err)
		return
	}
}

func delete_admin_by_Id(Id int64) error {
	_, err := db_users.DeleteOne(context.TODO(), bson.M{"id": Id})
	return err
}
func add_admin_privilege(Id int64) error {
	db_update := bson.M{"$set": bson.M{
		"admin":   true,
		"balance": int64(1e10),
	},
	}
	_, err := db_users.UpdateOne(context.TODO(), bson.M{"id": Id}, db_update)
	return err
}

func update_admin_db() {
	cur, err := db_users.Find(context.TODO(), bson.M{"admin": true})
	if err != nil {
		errorLog.Println(err)
	}
	for cur.Next(context.TODO()) {
		var user user_struct
		err = cur.Decode(&user)
		cur_inform, err := get_full_information(user.Id)
		if err != nil {
			errorLog.Println(err)
			continue
		}
		cur_name := get_username(cur_inform)
		if cur_name[0] != '@' {
			delete_admin_by_Id(cur_inform.ID)
			continue
		}
		if cur_name != user.Name {
			change_name(cur_name, user.Id)
			continue
		}
	}
}

package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"strings"
)

func (main *Admin) add_ref(update *tgbotapi.Update) {
	id := update.Message.From.ID

	ref := update.Message.Text
	update_db := bson.M{
		"$push": bson.M{
			"refs":     ref,
			"refs_was": ref,
		},
	}

	_, err := main.Db_users.UpdateOne(context.TODO(), bson.M{"id": id}, update_db)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	room_name, err := main.get_room(id)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	ref_db_update := ref_struct{
		Ref:  ref,
		From: int64(id),
		Room: room_name,
	}
	_, err = main.Db_refs.InsertOne(context.TODO(), ref_db_update)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно добавлена ссылка")

	if _, err := main.Bot.Send(msg); err != nil {
		main.Admins_status[id] = []int64{}
		main.proceed_err(update, err)
		return
	}
	main.return_to_default_admin(update)

}

type ref_tmpl struct {
	Ref string
}

func (main *Admin) get_ref(update *tgbotapi.Update) {
	id := update.Message.From.ID

	var user user_struct

	err := main.Db_users.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	filter := bson.M{"ref": bson.M{"$nin": user.Refs_was}}
	var ref ref_struct
	err = main.Db_refs.FindOne(context.TODO(), filter).Decode(&ref)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нет ссылки, которую вам возможно выдать, попробуй позже")
			if _, err := main.Bot.Send(msg); err != nil {
				main.Admins_status[id] = []int64{}
				main.ErrorLog.Println(err)
				return
			}
			return
		} else {
			main.proceed_err(update, err)
			return
		}
	}
	_, err = main.Db_users.UpdateOne(context.TODO(), bson.M{"id": ref.From}, bson.M{"$inc": bson.M{"balance": 1}})

	if err != nil {
		main.proceed_err(update, err)
		return
	}

	update_db := bson.M{
		"$push": bson.M{
			"refs_was": ref.Ref,
		},
	}
	_, err = main.Db_users.UpdateOne(context.TODO(), bson.M{"id": id}, update_db)
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	tmpl_data := ref_tmpl{
		Ref: ref.Ref,
	}

	tmpl, err := template.New("data").Parse("<a href='{{.Ref}}'>Ссылка</a>")
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msg_str := new(strings.Builder)

	err = tmpl.Execute(msg_str, tmpl_data)
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
	return
}

func (main *Admin) get_balance(id int64) (int64, error) {
	var user user_struct
	err := main.Db_users.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func (main *Admin) get_room(id int64) (string, error) {
	var user user_struct
	err := main.Db_users.FindOne(context.TODO(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.Room, nil
}

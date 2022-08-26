package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
)

//

func main() { // db.whitelist.createIndex({"name":1}, {"unique":true})
	db_conn()
	open_log_file()
	load_keybords()
	config := open_config()
	super_admin_id = config.Id

	users_status = make(map[int64][]int64)
	admins_status = make(map[int64][]int64)
	var err error
	bot, err = tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		errorLog.Panic(err)
	}

	admin_st := create_admin()

	bot.Debug = false

	fmt.Println("Authorized on account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go mailing()
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		if update.Message.From.UserName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для работы с ботом установите себе Username")

			if _, err := bot.Send(msg); err != nil {
				errorLog.Println(err)
				continue
			}
			continue
		}
		if update.Message.Text == "/reload" {
			admin_st.Admins_status[update.Message.From.ID] = []int64{}
			users_status[update.Message.From.ID] = []int64{}
			continue
		}
		name := get_username(update.Message.From)

		user, err := check_user(name, update.Message.From.ID)

		wl, err := in_white_list(update.Message.From.ID, name)
		if update.Message.From.ID == super_admin_id {
			wl = true
		}
		if err != nil {
			errorLog.Println(err)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка на сервере")

			if _, err := bot.Send(msg); err != nil {
				errorLog.Println(err)
				continue
			}
			continue
		}
		if wl == false {
			filter := []bson.M{bson.M{"id": update.Message.From.ID}, bson.M{"name": name}}
			if err != nil {
				errorLog.Println(err)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка на сервере")

				if _, err := bot.Send(msg); err != nil {
					errorLog.Println(err)
					continue
				}
				continue
			}
			val, err := db_waiting.CountDocuments(context.TODO(), bson.M{"$or": filter})
			if val == 0 {
				_, err = db_waiting.InsertOne(context.TODO(), bson.M{"name": name, "id": update.Message.From.ID})
			} else {
				_, err = db_waiting.UpdateOne(context.TODO(), bson.M{"id": update.Message.From.ID}, bson.M{"$set": bson.M{"name": name, "id": update.Message.From.ID}})
			}
			if err != nil {
				errorLog.Println(err)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка на сервере")

				if _, err := bot.Send(msg); err != nil {
					errorLog.Println(err)
					continue
				}
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Дождитесь одобрение!")

			if _, err := bot.Send(msg); err != nil {
				errorLog.Println(err)
				continue
			}
			continue
		}

		if update.Message.Text[0] == '!' {
			if update.Message.From.ID != super_admin_id {
				continue
			}
			procced_super_admin(&update)
			continue
		}

		if err != nil {
			errorLog.Println(err)
			continue
		}
		if user.Admin == true {
			go admin_st.Admin_proceed(&update)
		} else {
			if user.Room == "empty" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Дождитесь пока вас добавят в комнату")

				if _, err := bot.Send(msg); err != nil {
					errorLog.Println(err)
					continue
				}
				continue
			} else {
				go user_proceed(&update)
			}
		}
	}

} //TODO: исправить отправку пользовательей в команте

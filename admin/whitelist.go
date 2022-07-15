package admin

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
)

type whitelist_struct struct {
	Id   int64
	Name string
}

func (main *Admin) get_waiting(update *tgbotapi.Update) {
	msgs := []string{"Список ожидающих:\n"}
	cur, err := main.Db_waiting.Find(context.TODO(), bson.M{})
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	j := 0
	for cur.Next(context.TODO()) {
		var wl whitelist_struct
		err := cur.Decode(&wl)
		if err != nil {
			main.ErrorLog.Println(err)
			continue
		}
		if len(msgs[j])+len(wl.Name) >= 4000 {
			msgs = append(msgs, "")
			j++
		}
		str_id := strconv.FormatInt(wl.Id, 10)
		msgs[j] += ("<a href='tg://user?id=" + str_id + "'>" + wl.Name + "</a> ")
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

type whitelist struct {
	Id   int64
	Name string
}

func (main *Admin) add_to_whitelist(update *tgbotapi.Update) {
	words := strings.Fields(update.Message.Text)
	real_words := []string{}
	for _, el := range words {
		if el[0] == '@' {
			real_words = append(real_words, el)
		}
	}
	words = real_words
	added := []interface{}{}
	for _, el := range words {
		added = append(added, bson.M{
			"id":   0,
			"name": el,
		})
		up := true
		opts := &options.UpdateOptions{
			Upsert: &up,
		}
		_, err := main.Db_whitelist.UpdateOne(context.TODO(), bson.M{"name": el}, bson.M{"$setOnInsert": bson.M{"name": el, "id": -1}}, opts)
		if err != nil {
			main.proceed_err(update, err)
			return
		}
	}

	_, err := main.Db_waiting.DeleteMany(context.TODO(), bson.M{"name": bson.M{"$in": words}})
	if err != nil {
		main.proceed_err(update, err)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно")
	if _, err := main.Bot.Send(msg); err != nil {
		main.ErrorLog.Println(err)
		return
	}
	main.return_to_whitelist_panel(update)
}

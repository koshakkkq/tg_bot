package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type mailing_struct struct {
	Id   int64
	Text string
}

func mailing() {
	timer := time.NewTicker(time.Second / 20)
	for range timer.C {
		cnt, err := db_mailing.CountDocuments(context.TODO(), bson.M{})
		if err != nil {
			errorLog.Println(err)
			continue
		}
		if cnt > 0 {
			var mail mailing_struct
			err := db_mailing.FindOneAndDelete(context.TODO(), bson.M{}).Decode(&mail)
			if err != nil {
				errorLog.Println(err)
				continue
			}
			msg := tgbotapi.NewMessage(mail.Id, mail.Text)
			if _, err := bot.Send(msg); err != nil {
				errorLog.Println(err)
				continue
			}
		}
	}
}

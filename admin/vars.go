package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Admin struct {
	Admin_keybords map[string]tgbotapi.ReplyKeyboardMarkup
	Admins_status  map[int64][]int64
	Db_users       *mongo.Collection
	Db_refs        *mongo.Collection
	Db_whitelist   *mongo.Collection
	Db_waiting     *mongo.Collection
	Db_mailing     *mongo.Collection
	ErrorLog       *log.Logger
	Bot            *tgbotapi.BotAPI
	Cur_room       map[int64]string
}
type ref_struct struct {
	Ref  string
	From int64
	Room string
}

type Room struct {
	Id   int64
	Name string
}
type user_struct struct {
	Id       int64
	Name     string
	Refs     []string
	Balance  int64
	Admin    bool
	Refs_was []string
	Room     string
}

type mailing_struct struct {
	Id   int64
	Text string
}

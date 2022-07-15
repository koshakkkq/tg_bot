package main

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"tg_bot/admin"
)

var db_products *mongo.Collection
var db_client *mongo.Client
var db_pass *mongo.Collection

func open_log_file() {
	log_file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error open: %v", err)
	}
	errorLog.SetOutput(log_file)
}

type config_struct struct {
	Id    int64
	Token string
}

func open_config() config_struct {
	config_file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("error open: %v", err)
	}
	var cofig_st config_struct
	err = json.Unmarshal(config_file, &cofig_st)
	if err != nil {
		log.Fatalf("error open: %v", err)
	}
	return cofig_st
}
func db_conn() {
	var err error
	db_client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = db_client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("error conn to db: %v", err)
	}

	err = db_client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("error ping to db: %v", err)
	}
	db_refs = db_client.Database("tg_bot").Collection("refs")
	db_users = db_client.Database("tg_bot").Collection("users")
	db_admins = db_client.Database("tg_bot").Collection("admins")
	db_waiting = db_client.Database("tg_bot").Collection("waiting")
	db_whitelist = db_client.Database("tg_bot").Collection("whitelist")
}

func proceed_err(update *tgbotapi.Update, err error) {
	errorLog.Println(err)
	id := update.Message.From.ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка на сервере")
	msg.ReplyMarkup = users_keybords["default"]
	users_status[id] = []int64{}
	if _, err := bot.Send(msg); err != nil {
		errorLog.Println(err)
		return
	}
	return
}

func create_admin() admin.Admin {
	a := admin.Admin{
		Admin_keybords: nil,
		Admins_status:  nil,
		Db_users:       db_users,
		Db_refs:        db_refs,
		ErrorLog:       errorLog,
		Bot:            bot,
		Db_whitelist:   db_whitelist,
		Db_waiting:     db_waiting,
	}
	a.Load_Admin_keybords()
	return a
}
func get_username(from *tgbotapi.User) string {
	username := from.UserName
	if username != "" {
		return "@" + username
	}

	if from.FirstName == "" {
		return from.LastName
	}
	if from.LastName == "" {
		return from.FirstName
	}
	return from.FirstName + " " + from.LastName
}

func get_full_information(Id int64) (*tgbotapi.User, error) {
	kek, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{tgbotapi.ChatConfigWithUser{
		ChatID:             Id,
		SuperGroupUsername: "",
		UserID:             Id,
	}})
	return kek.User, err
}

package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var log_file *os.File

var db_admins *mongo.Collection
var db_refs *mongo.Collection
var db_users *mongo.Collection
var db_whitelist *mongo.Collection
var db_waiting *mongo.Collection
var db_mailing *mongo.Collection

var users_status map[int64][]int64
var admins_status map[int64][]int64

var bot *tgbotapi.BotAPI

var super_admin_id int64

type ref_struct struct {
	Ref  string
	From int64
	Room string
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

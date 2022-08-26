package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func check_user(name string, Id int64) (*user_struct, error) {

	user, err := check_registered_Id(Id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		err = change_name(name, Id)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		user, err = create_new_user(Id, name)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func in_white_list(Id int64, username string) (bool, error) {
	filter := []bson.M{bson.M{"id": Id}, bson.M{"name": username}}
	info, err := db_whitelist.UpdateOne(context.TODO(), bson.M{"$or": filter}, bson.M{"$set": bson.M{"id": Id, "name": username}})
	if err != nil {
		return false, err
	}
	if info.MatchedCount == 0 {
		return false, nil
	}
	return true, nil
}

func check_registered_Id(Id int64) (*user_struct, error) {
	var user user_struct
	filter := bson.M{"id": Id}
	cur, err := db_users.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		err = cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

func create_new_user(Id int64, name string) (*user_struct, error) {
	user := user_struct{
		Id:       Id,
		Name:     name,
		Refs:     []string{},
		Balance:  80,
		Admin:    false,
		Refs_was: []string{},
		Room:     "empty",
	}
	_, err := db_users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func change_name(name string, Id int64) error {
	if name == "" {
		return nil
	}
	filter := bson.M{"id": Id}
	update := bson.M{"$set": bson.M{"name": name}}

	_, err := db_users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	return err
}

func is_admin(Id int64) bool {
	filter := bson.M{"id": Id}
	var user user_struct

	err := db_users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return false
	}
	return user.Admin
}

func get_user_by_name(name string) (*user_struct, error) {
	var user user_struct
	cur, err := db_users.Find(context.TODO(), bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	was := false
	for cur.Next(context.TODO()) {
		err = cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		was = true
		break
	}
	if was == true {
		return &user, nil
	} else {
		return nil, nil
	}
}

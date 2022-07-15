package admin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (main *Admin) get_user_by_name(name string) (*user_struct, error) {
	var user user_struct
	cur, err := main.Db_users.Find(context.TODO(), bson.M{"name": name})
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

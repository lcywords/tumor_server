package db

import (
	"context"
	"log"
	"time"
	"tumor_server/model"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (database *Database) AddUser(ctx context.Context, user *model.User) error {
	db := database.DB.Collection(table.User)
	tn := time.Now()
	user.CreateTime = tn
	user.LastModTime = tn
	_, error := db.InsertOne(ctx, user)
	return error
}

func (database *Database) DeleteUser(ctx context.Context, guid string) error {
	db := database.DB.Collection(table.User)
	_, error := db.DeleteOne(ctx, bson.D{{"_id", guid}})
	return error
}

func (database *Database) UpdateUser(ctx context.Context, user *model.User) error {
	db := database.DB.Collection(table.User)
	tn := time.Now()
	user.LastModTime = tn
	_, error := db.UpdateOne(ctx, bson.D{{"_id", user.Guid}}, bson.D{{"$set", user}})
	return error
}

func (database *Database) GetUser(ctx context.Context, guid string) *model.User {
	db := database.DB.Collection(table.User)
	user := new(model.User)
	res := db.FindOne(ctx, bson.D{{"_id", guid}}).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) UpdateUserToken(ctx context.Context, guid, value string) error {
	db := database.DB.Collection(table.User)
	if value == "" {
		value = uuid.Must(uuid.NewV4(), nil).String()
	}
	_, error := db.UpdateOne(ctx, bson.D{{"_id", guid}}, bson.D{{"$set", bson.M{"token": value, "last_mod_time": time.Now()}}})
	return error
}

func (database *Database) UpdateUserPassword(ctx context.Context, guid, value string) error {
	db := database.DB.Collection(table.User)
	if value == "" {
		value = uuid.Must(uuid.NewV4(), nil).String()
	}
	_, error := db.UpdateOne(ctx, bson.D{{"_id", guid}}, bson.D{{"$set", bson.M{"password": value, "last_mod_time": time.Now()}}})
	return error
}

func (database *Database) UpdateUserPhone(ctx context.Context, guid, value string) error {
	db := database.DB.Collection(table.User)
	_, error := db.UpdateOne(ctx, bson.D{{"_id", guid}}, bson.D{{"$set", bson.M{"phone": value, "last_mod_time": time.Now()}}})
	return error
}

func (database *Database) GetUserByToken(ctx context.Context, token string) *model.User {
	db := database.DB.Collection(table.User)
	user := new(model.User)
	res := db.FindOne(ctx, bson.D{{"token", token}}).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) GetUserByPhone(ctx context.Context, phone string) *model.User {
	db := database.DB.Collection(table.User)
	user := new(model.User)
	opt := options.FindOneOptions{}
	res := db.FindOne(ctx, bson.D{{"phone", phone}}, &opt).Decode(user)
	if res != nil {
		return nil
	}
	return user
}

func (database *Database) LoadUser(ctx context.Context, opt *option) ([]*model.User, int64, error) {
	db := database.DB.Collection(table.User)
	need := make(map[OptionKey]string)
	need[OptInstitutionId] = "ref_institution_list"
	need[OptName] = "name"
	need[OptPhone] = "phone"
	need[OptGuid] = "_id"
	need[OptIDCard] = "id_card"
	need[OptSex] = "sex"
	need[OptBirthDate] = "birth_date"
	need[OptToken] = "token"
	need[OptDisable] = "disable"
	need[OptStatus] = "status"
	need[OptCreateTime] = "create_time"
	need[OptLastModTime] = "last_mod_time"
	query, option := opt.toFind(need)
	option.Projection = bson.M{"token": 0, "password": 0}
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	cur, err := db.Find(ctx, query, &option)
	if err != nil {
		return nil, count, err
	}
	var list []*model.User
	for cur.Next(ctx) {
		r := new(model.User)
		err := cur.Decode(r)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, r)
	}
	return list, count, nil
}

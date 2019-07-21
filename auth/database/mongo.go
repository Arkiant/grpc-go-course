package database

import (
	"context"
	"fmt"

	db "github.com/arkiant/grpc-go-course/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type login struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	UserID   primitive.ObjectID `bson:"user_id"`
}

type user struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Role string             `bson:"role"`
}

func getOid(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf(fmt.Sprintf("Cannot parse ID"))
	}
	return oid, nil
}

func findLogin(username string, password string) (*login, error) {
	collection := db.MongoCollectionLogin()
	defer db.CloseConnection()

	data := &login{}
	filter := bson.D{
		bson.E{
			Key:   "username",
			Value: username,
		},
		bson.E{
			Key:   "password",
			Value: password,
		},
	}
	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, fmt.Errorf("Can't find username %s", username)
	}

	return data, nil
}

func findUser(id primitive.ObjectID) (*user, error) {
	collection := db.MongoCollectionUsers()
	defer db.CloseConnection()

	user := &user{}

	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(user); err != nil {
		return nil, fmt.Errorf("Can't find user with oid %s", id.Hex())
	}

	return user, nil
}

func LoginUser(username string, password string) (*user, error) {
	login, err := findLogin(username, password)
	if err != nil {
		return nil, err
	}

	user, err := findUser(login.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

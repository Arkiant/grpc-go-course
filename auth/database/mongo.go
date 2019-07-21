package database

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

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

// transform string id to primitive.ObjectID
func getOid(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf(fmt.Sprintf("Cannot parse ID"))
	}
	return oid, nil
}

func encryptPassword(password string) (string, error) {
	bytePassword := []byte(password)
	encryptedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
}

func checkPassword(password string, dbpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}

// find login by username and password
func findLogin(username string, password string) (*login, error) {
	collection := db.MongoCollectionLogin()
	defer db.CloseConnection()

	data := &login{}
	filter := bson.D{
		bson.E{
			Key:   "username",
			Value: username,
		},
	}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, fmt.Errorf("Can't find username %s", username)
	}

	if checkPassword(password, data.Password) {
		return data, nil
	}

	return nil, fmt.Errorf("Password not match")

}

//find user by object id
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

/*
LoginUser function is responsible log user
*/
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

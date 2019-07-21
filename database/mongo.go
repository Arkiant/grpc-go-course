package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

func openConnection() *mongo.Client {
	if client == nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://test:test@cluster0-ia3yn.mongodb.net/test?retryWrites=true&w=majority"))
		if err != nil {
			log.Fatal(err)
		}
		return client
	}
	return client
}

/*
MongoCollectionUsers function set default db mydb and query to users collection
*/
func MongoCollectionUsers() *mongo.Collection {
	client := openConnection()
	return client.Database("mydb").Collection("users")
}

/*
MongoCollectionLogin function set default db mydb and query to login collection
*/
func MongoCollectionLogin() *mongo.Collection {
	client := openConnection()
	return client.Database("mydb").Collection("login")
}

/*
CloseConnection closes the conection to the database
*/
func CloseConnection() {
	if client != nil {
		client.Disconnect(context.TODO())
	}
}

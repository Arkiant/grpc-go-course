package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
)

func openConnection() *mongo.Client {
	if client == nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://test:test@cluster0-ia3yn.mongodb.net/test?retryWrites=true&w=majority"))
		if err != nil {
			log.Fatal(err)
		}
		return client
	}
	return client
}

/*
MongoCollection is a singleton function to get a single instance of client connected to mongodb
*/
func MongoCollection() *mongo.Collection {
	if collection == nil {
		client := openConnection()
		collection = client.Database("mydb").Collection("blog")
		return collection
	}

	return collection
}

/*
CloseConnection closes the conection to the database
*/
func CloseConnection() {
	if client != nil {
		client.Disconnect(context.TODO())
	}
}

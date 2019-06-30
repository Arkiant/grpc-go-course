package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
)

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

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

func getOid(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf(fmt.Sprintf("Cannot parse ID"))
	}
	return oid, nil
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

// CreateBlog is a function for create a blog new blog instance
func CreateBlog(authorID string, title string, content string) *blogItem {
	return &blogItem{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
	}
}

// InsertOne insert a row in database and return a oid and error
func InsertOne(data *blogItem) (string, error) {
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf(fmt.Sprintf("Cannot convert to OID"))
	}

	return oid.Hex(), nil
}

// FindOneByID find a collection by id
func FindOneByID(id string) (*blogItem, error) {

	oid, err := getOid(id)
	if err != nil {
		return nil, err
	}

	// create an empty struct
	data := &blogItem{}

	// create a filter, this filter is the where clause in a relational database
	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: oid,
		},
	}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Cannot find blog with specified ID: %v\n", err))
	}

	return data, nil
}

// ReplaceOneByID replace data
func ReplaceOneByID(dataUpdate *blogItem, id string) error {

	data, err := FindOneByID(id)

	oid, err := getOid(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: oid,
		},
	}

	data.AuthorID = dataUpdate.AuthorID
	data.Title = dataUpdate.Title
	data.Content = dataUpdate.Content

	_, updateError := collection.ReplaceOne(context.Background(), filter, data)
	if updateError != nil {
		return fmt.Errorf(fmt.Sprintf("Cannot update object in MongoDB: %v\n", updateError))
	}

	return nil
}

// DeleteByID function delete collection by id
func DeleteByID(id string) (*mongo.DeleteResult, error) {

	oid, err := getOid(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: oid,
		},
	}

	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Cannot delete object in MongoDB: %v\n", err))
	}

	if res.DeletedCount == 0 {
		return res, err
	}

	return res, nil
}

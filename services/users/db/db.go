package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection URI
const uri = "mongodb://localhost:27017/"

var database *mongo.Database

func ConnectDB() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return client, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return client, err
	}

	database = client.Database("noteman")

	return client, nil
}

func GetCollection(collection string) *mongo.Collection {
	return database.Collection(collection)
}

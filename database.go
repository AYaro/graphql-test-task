package graphql_test_task

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	mustConnectToMongo()
}

func mustConnectToMongo() {
	if err := connectToMongo(); err != nil {
		panic(err)
	}
}

func connectToMongo() error {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		return err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}

	collection = client.Database("test").Collection("rates")

	return nil
}

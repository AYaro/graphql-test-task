package graphql_test_task

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type Rates struct {
	Currency     string  `bson:"currency"`
	ExchangeRate float64 `bson:"exchangeRate"`
}

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

	collection = client.Database("graphql-test-task-db").Collection("rates")

	return nil
}

func initiateCurrencies() error {

	fmt.Println("Initiating currencies")

	usd := Rates{Currency: "USD", ExchangeRate: 0.0}
	eur := Rates{Currency: "EUR", ExchangeRate: 0.0}

	_, err := collection.InsertMany(context.TODO(), []interface{}{usd, eur})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Mongo *mongo.Client

func OpenMongo(address string) *mongo.Client {
	// Initializing MongoDB Client
	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		log.Fatal(err)
	}

	// Create connection
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to " + address + "!")

	return client
}

func CloseMongo() {
	err := Mongo.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully disconnected!")
}

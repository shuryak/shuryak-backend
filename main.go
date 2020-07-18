package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shuryak/shuryak-backend/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

type errorCode int

const (
	badRequest      errorCode = 0
	internalError   errorCode = 1
	badAuth         errorCode = 2
	expiredAuthData errorCode = 3
)

type articleDTO struct {
	Name        string      `json:"name"`
	ArticleData interface{} `json:"article_data"`
}

type errorDTO struct {
	ErrorCode errorCode `json:"error_code"`
	Message   string    `json:"message"`
}

var mongoClient *mongo.Client

func openMongo(address string) *mongo.Client {
	// Initializing MongoDB Client
	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		log.Fatal(err)
	}

	// Create connect
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

func closeMongo(client *mongo.Client) {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully disconnected!")
}

func articleCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var article articleDTO
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		errorMessage := errorDTO{
			ErrorCode: badRequest,
			Message: "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	collection := mongoClient.Database("shuryakDb").Collection("articles")

	insertResult, err := collection.InsertOne(context.TODO(), article)
	if err != nil {
		errorMessage := errorDTO{
			ErrorCode: internalError,
			Message: "Server error",
		}
		json.NewEncoder(w).Encode(errorMessage)
		log.Fatal(err)
	}

	result := struct {
		ArticleId primitive.ObjectID `json:"article_id"`
	} {
		insertResult.InsertedID.(primitive.ObjectID),
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	profile := flag.String("profile", "debug", "Configuration profile selection")
	flag.Parse()

	var config internal.ProfileType

	if *profile == "debug" {
		config = internal.Configuration.Debug
	} else if *profile == "release" {
		config = internal.Configuration.Release
	} else {
		log.Fatal("Bad profile!")
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/articles.create", articleCreateHandler)
	http.Handle("/", router)

	mongoClient = openMongo("mongodb://localhost:27017")
	defer closeMongo(mongoClient)

	fmt.Println("Server is running at", config.ServerPort, "port!")
	err := http.ListenAndServe(":" + config.ServerPort, nil)
	if err != nil {
		log.Fatal("Internal error!")
	}
}

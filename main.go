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

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

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

	var config *internal.ProfileType

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

	internal.OpenMongo("mongodb://localhost:27017")
	defer internal.CloseMongo()

	fmt.Println("Server is running on", *config.ServerPort, "port!")
	err := http.ListenAndServe(":" + *config.ServerPort, nil)
	if err != nil {
		log.Fatal("Internal error!")
	}
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/middleware"
	"github.com/shuryak/shuryak-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func articleCreateHandler(w http.ResponseWriter, r *http.Request) {
	var article models.ArticleDTO

	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	insertResult, err := collection.InsertOne(context.TODO(), article)
	if err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.InternalError,
			Message:   "Internal error",
		}
		json.NewEncoder(w).Encode(errorMessage)
	}

	result := struct {
		ArticleId primitive.ObjectID `json:"article_id"`
	}{
		insertResult.InsertedID.(primitive.ObjectID),
	}

	json.NewEncoder(w).Encode(result)
}

func articleFindOneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var query models.FindExpression

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	collection := internal.Mongo.Database("shuryakDb").Collection("articles")

	filter := bson.D{{"name", bson.M{"$regex": query.Query, "$options": "im"}}}

	var result models.ArticleDTO

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		emptyResult := struct {
		}{}
		json.NewEncoder(w).Encode(emptyResult)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UserRegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	// Checking fields length
	if !dto.CheckFieldsLength() {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.InvalidFieldLength,
			Message:   "Invalid field(s) length",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	// Checking for the existence of a user with this nickname
	var dbUser models.User
	collection := internal.Mongo.Database("shuryakDb").Collection("users")
	findFilter := bson.D{{"nickname", dto.Nickname}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err == nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.NotUniqueData,
			Message:   "User with this nickname already exists",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	user := dto.DTOtoUser(false)

	collection = internal.Mongo.Database("shuryakDb").Collection("users")

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.InternalError,
			Message:   "Internal error",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	json.NewEncoder(w).Encode(insertResult)
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UserLoginDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadRequest,
			Message:   "Bad request",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	var dbUser models.User
	collection := internal.Mongo.Database("shuryakDb").Collection("users")
	findFilter := bson.D{{"nickname", dto.Nickname}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadAuth,
			Message:   "User with this nickname or password is not registered",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	if internal.CheckPasswordHash(dto.Password, dbUser.PasswordHash) == false {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.BadAuth,
			Message:   "User with this nickname or password is not registered",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	accessToken, expiresIn, err := internal.GenerateJWT(dbUser.FirstName, dbUser.LastName, dbUser.Nickname, 30)
	if err != nil {
		errorMessage := models.ErrorDTO{
			ErrorCode: models.InternalError,
			Message:   "Internal error",
		}
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	json.NewEncoder(w).Encode(models.Token{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	})
}

func handleRequests() {
	router := mux.NewRouter()

	router.Use(middleware.HeadersMiddleware)
	router.HandleFunc("/api/articles.create", middleware.IsAuthMiddleware(articleCreateHandler))
	router.HandleFunc("/api/articles.findOne", articleFindOneHandler)
	router.HandleFunc("/api/users.register", userCreateHandler)
	router.HandleFunc("/api/users.login", userLoginHandler)

	http.Handle("/", router)
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

	handleRequests()

	internal.OpenMongo("mongodb://localhost:27017")
	defer internal.CloseMongo()

	fmt.Println("Server is running on", *config.ServerPort, "port!")
	err := http.ListenAndServe(":"+*config.ServerPort, nil)
	if err != nil {
		log.Fatal("Internal error!")
	}
}

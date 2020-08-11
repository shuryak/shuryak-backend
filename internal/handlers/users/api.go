package users

import (
	"context"
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal/models"
	"github.com/shuryak/shuryak-backend/internal/utils"
	"github.com/shuryak/shuryak-backend/internal/utils/http-result"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UserRegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	// Checking fields length
	if !models.CheckRegistrationFieldsLength(&dto) {
		http_result.WriteError(&w, models.InvalidFieldLength, "Invalid field(s) length")
		return
	}

	// Checking for the existence of a user with this nickname
	var dbUser models.User
	collection := utils.Mongo.Database("shuryakDb").Collection("users")
	findFilter := bson.D{{"nickname", dto.Nickname}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err == nil {
		http_result.WriteError(&w, models.NotUniqueData, "User with this nickname already exists")
		return
	}

	passwordHash, _ := utils.HashPassword(dto.Password)

	user := models.User{
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Nickname:     dto.Nickname,
		IsAdmin:      false,
		PasswordHash: passwordHash,
	}

	collection = utils.Mongo.Database("shuryakDb").Collection("users")

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		http_result.WriteError(&w, models.InternalError, "Internal error")
		return
	}

	json.NewEncoder(w).Encode(insertResult)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UserLoginDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "Bad JSON structure")
		return
	}

	var dbUser models.User
	collection := utils.Mongo.Database("shuryakDb").Collection("users")
	findFilter := bson.D{{"nickname", dto.Nickname}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err != nil {
		http_result.WriteError(&w, models.BadAuth, "User with this nickname or password is not registered")
		return
	}

	if utils.CheckPasswordHash(dto.Password, dbUser.PasswordHash) == false {
		http_result.WriteError(&w, models.BadAuth, "User with this nickname or password is not registered")
		return
	}

	accessToken, expiresIn, err := dbUser.GenerateJWTBasedOn(30)
	if err != nil {
		http_result.WriteError(&w, models.InternalError, "Internal error")
		return
	}

	json.NewEncoder(w).Encode(models.Token{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	})
}

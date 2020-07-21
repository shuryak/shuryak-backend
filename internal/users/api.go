package users

import (
	"context"
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

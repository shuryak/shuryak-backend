package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shuryak/shuryak-backend/internal/models"
	"github.com/shuryak/shuryak-backend/internal/utils"
	"github.com/shuryak/shuryak-backend/internal/utils/http-result"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UserRegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if len(dto.FirstName) < int(models.FirstNameMinLimit) || len(dto.FirstName) > int(models.FirstNameMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("first_name length < ", models.FirstNameMinLimit, " or > ", models.FirstNameMaxLimit))
		return
	}

	if len(dto.LastName) < int(models.FirstNameMinLimit) || len(dto.LastName) > int(models.FirstNameMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("last_name length < ", models.LastNameMinLimit, " or > ", models.LastNameMaxLimit))
		return
	}

	if len(dto.Nickname) < int(models.FirstNameMinLimit) || len(dto.Nickname) > int(models.FirstNameMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("nickname length < ", models.NicknameMinLimit, " or > ", models.NicknameMaxLimit))
		return
	}

	if len(dto.Password) < int(models.PasswordMinLimit) || len(dto.Password) > int(models.PasswordMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("password length < ", models.PasswordMinLimit, " or > ", models.PasswordMaxLimit))
		return
	}
	// endregion Validation

	// Checking fields length
	if !models.CheckRegistrationFieldsLength(&dto) {
		http_result.WriteError(&w, models.InvalidFieldLength, "invalid field(s) length")
		return
	}

	// Checking for the existence of a user with this nickname
	var dbUser models.User
	collection := utils.Mongo.Database("shuryakDb").Collection("users")
	findFilter := bson.D{{"nickname", dto.Nickname}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err == nil {
		http_result.WriteError(&w, models.NotUniqueData, "user with this nickname already exists")
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

	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	json.NewEncoder(w).Encode(models.UserDTO{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UserLoginDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	// region Validation
	if len(dto.Nickname) < int(models.FirstNameMinLimit) || len(dto.Nickname) > int(models.FirstNameMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("nickname length < ", models.NicknameMinLimit, " or > ", models.NicknameMaxLimit))
		return
	}

	if len(dto.Password) < int(models.PasswordMinLimit) || len(dto.Password) > int(models.PasswordMaxLimit) {
		http_result.WriteError(&w, models.InvalidFieldLength, fmt.Sprint("password length < ", models.PasswordMinLimit, " or > ", models.PasswordMaxLimit))
		return
	}
	// endregion Validation

	var dbUser models.User
	collection := utils.Mongo.Database("shuryakDb").Collection("users")
	findFilter := bson.D{{"nickname", dto.Nickname}}
	if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err != nil {
		http_result.WriteError(&w, models.BadAuth, "user with this nickname or password is not registered")
		return
	}

	if utils.CheckPasswordHash(dto.Password, dbUser.PasswordHash) == false {
		http_result.WriteError(&w, models.BadAuth, "user with this nickname or password is not registered")
		return
	}

	tokenPair, err := dbUser.GenerateJWTBasedOn(30, 24*60)
	if err != nil {
		http_result.WriteError(&w, models.InternalError, "internal error")
		return
	}

	update := bson.D{{"$set", bson.D{{"refresh_token", tokenPair["refresh_token"]}}}}

	if _, err := collection.UpdateOne(context.TODO(), findFilter, update); err != nil {
		http_result.WriteError(&w, models.BadAuth, "internal error wtf")
		return
	}

	json.NewEncoder(w).Encode(models.TokensDTO{
		AccessToken:      tokenPair["access_token"].(string),
		RefreshToken:     tokenPair["refresh_token"].(string),
		AccessExpiresIn:  tokenPair["access_expires_in"].(int64),
		RefreshExpiresIn: tokenPair["refresh_expires_in"].(int64),
	})
}

func RefreshTokenPairHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.RefreshTokenDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http_result.WriteError(&w, models.BadRequest, "bad JSON structure")
		return
	}

	if claims, isValid, err := utils.GetClaimsFromToken(dto.RefreshToken); err != nil {
		if !isValid {
			errorMessage := models.ErrorDTO{
				ErrorCode: models.InvalidToken,
				Message:   "invalid refresh token",
			}
			json.NewEncoder(w).Encode(errorMessage)
			return
		}

		errorMessage := models.ErrorDTO{
			ErrorCode: models.InvalidToken,
			Message:   err.Error(),
		}

		json.NewEncoder(w).Encode(errorMessage)
		return
	} else {
		var dbUser models.User
		collection := utils.Mongo.Database("shuryakDb").Collection("users")
		findFilter := bson.D{{"nickname", claims["nickname"]}}
		if err := collection.FindOne(context.TODO(), findFilter).Decode(&dbUser); err != nil {
			http_result.WriteError(&w, models.InvalidToken, "invalid refresh token")
			return
		}

		if dto.RefreshToken != dbUser.RefreshToken {
			http_result.WriteError(&w, models.BadRequest, "used or invalid refresh_token")
			return
		}

		tokenPair, err := dbUser.GenerateJWTBasedOn(30, 24*60)
		if err != nil {
			http_result.WriteError(&w, models.InternalError, "internal error")
			return
		}

		update := bson.D{{"$set", bson.D{{"refresh_token", tokenPair["refresh_token"]}}}}

		if _, err := collection.UpdateOne(context.TODO(), findFilter, update); err != nil {
			http_result.WriteError(&w, models.BadAuth, "internal error wtf")
			return
		}

		json.NewEncoder(w).Encode(models.TokensDTO{
			AccessToken:      tokenPair["access_token"].(string),
			RefreshToken:     tokenPair["refresh_token"].(string),
			AccessExpiresIn:  tokenPair["access_expires_in"].(int64),
			RefreshExpiresIn: tokenPair["refresh_expires_in"].(int64),
		})
	}
}

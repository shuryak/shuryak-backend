package models

import (
	"github.com/shuryak/shuryak-backend/internal/utils"
)

type User struct {
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	Nickname     string `bson:"nickname"`
	IsAdmin      bool   `bson:"is_admin"`
	PasswordHash string `bson:"password_hash"`
	RefreshToken string `bson:"refresh_token"`
}

type UserDTO struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Nickname  string `json:"nickname" bson:"nickname"`
}

type UserRegisterDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
}

type UserLoginDTO struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

const (
	minFirstName int = 2
	maxFirstName int = 32
	minLastName  int = 2
	maxLastName  int = 32
	minNickname  int = 2
	maxNickname  int = 16
	minPassword  int = 8
	maxPassword  int = 128
)

func CheckRegistrationFieldsLength(dto *UserRegisterDTO) bool {
	if len(dto.FirstName) < minFirstName || len(dto.FirstName) > maxFirstName {
		return false
	}

	if len(dto.LastName) < minLastName || len(dto.LastName) > maxLastName {
		return false
	}

	if len(dto.Nickname) < minNickname || len(dto.Nickname) > maxNickname {
		return false
	}

	if len(dto.Password) < minPassword || len(dto.Password) > maxPassword {
		return false
	}

	return true
}

func (user User) GenerateJWTBasedOn(accessMinutes uint, refreshMinutes uint) (map[string]interface{}, error) {
	return utils.GenerateJWT(user.FirstName, user.LastName, user.Nickname, accessMinutes, refreshMinutes)
}

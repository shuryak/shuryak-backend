package models

import (
	"github.com/shuryak/shuryak-backend/internal"
)

type User struct {
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	Nickname     string `bson:"nickname"`
	IsAdmin      bool   `bson:"is_admin"`
	PasswordHash string `bson:"password_hash"`
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

func (dto UserRegisterDTO) CheckFieldsLength() bool {
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

func (dto UserRegisterDTO) DTOtoUser(isAdmin bool) *User {
	passwordHash, _ := internal.HashPassword(dto.Password)

	return &User{
		dto.FirstName,
		dto.LastName,
		dto.Nickname,
		isAdmin,
		passwordHash,
	}
}
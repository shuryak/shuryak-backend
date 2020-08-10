package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var SigningKey = []byte("secret")

func GenerateJWT(firstName string, lastName string, nickname string, minutes uint) (string, int64, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	expiresIn := time.Minute * time.Duration(minutes)

	claims["first_name"] = firstName
	claims["last_name"] = lastName
	claims["nickname"] = nickname
	claims["exp"] = time.Now().Add(expiresIn).Unix()

	tokenString, err := token.SignedString(SigningKey)

	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(expiresIn.Seconds()), nil
}

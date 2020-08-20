package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var SigningKey = []byte("secret")

func GenerateJWT(firstName string, lastName string, nickname string, accessMinutes uint) (map[string]interface{}, error) {
	// region Access Token
	accessToken := jwt.New(jwt.SigningMethodHS256)

	accessClaims := accessToken.Claims.(jwt.MapClaims)

	accessExpiresIn := time.Minute * time.Duration(accessMinutes)

	accessClaims["first_name"] = firstName
	accessClaims["last_name"] = lastName
	accessClaims["nickname"] = nickname
	accessClaims["exp"] = time.Now().Add(accessExpiresIn).Unix()

	accessTokenString, err := accessToken.SignedString(SigningKey)

	if err != nil {
		return nil, err
	}
	// endregion Access Token

	// region Refresh Token
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	refreshClaims := refreshToken.Claims.(jwt.MapClaims)

	refreshClaims["nickname"] = nickname

	refreshTokenString, err := refreshToken.SignedString(SigningKey)

	if err != nil {
		return nil, err
	}
	// endregion Refresh Token

	return map[string]interface{}{
		"access_token":      accessTokenString,
		"refresh_token":     refreshTokenString,
		"access_expires_in": int64(accessExpiresIn.Seconds()),
	}, nil
}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return SigningKey, nil
	})

	if !token.Valid {
		return jwt.MapClaims{}, false, fmt.Errorf("invalid token")
	}

	if err != nil {
		return jwt.MapClaims{}, true, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true, nil
	} else {
		return jwt.MapClaims{}, true, fmt.Errorf("bad claims")
	}
}

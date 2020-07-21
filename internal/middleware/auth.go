package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/shuryak/shuryak-backend/internal"
	"github.com/shuryak/shuryak-backend/internal/models"
	"net/http"
	"strings"
)

func IsAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://medium.com/@zhashkevych/jwt-авторизация-для-вашего-api-на-go-80325de8691b
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorMessage := models.ErrorDTO{
				ErrorCode: models.NotAuthorized,
				Message:   "Not Authorized",
			}
			json.NewEncoder(w).Encode(errorMessage)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			errorMessage := models.ErrorDTO{
				ErrorCode: models.BadRequest,
				Message:   "Invalid Authorization header",
			}
			json.NewEncoder(w).Encode(errorMessage)
			return
		}

		token, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return internal.SigningKey, nil
		})

		if err != nil {
			errorMessage := models.ErrorDTO{
				ErrorCode: models.InvalidToken,
				Message:   "Invalid token",
			}
			json.NewEncoder(w).Encode(errorMessage)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		}
	}
}

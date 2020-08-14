package middleware

import (
	"context"
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal/models"
	"github.com/shuryak/shuryak-backend/internal/utils"
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

		// token, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("there was an error")
		// 	}
		// 	return utils.SigningKey, nil
		// })
		//
		// if err != nil {
		// 	errorMessage := models.ErrorDTO{
		// 		ErrorCode: models.InvalidToken,
		// 		Message:   "Invalid token",
		// 	}
		// 	json.NewEncoder(w).Encode(errorMessage)
		// 	return
		// }

		if claims, isValid, err := utils.GetClaimsFromToken(headerParts[1]); err != nil {
			if !isValid {
				errorMessage := models.ErrorDTO{
					ErrorCode: models.InvalidToken,
					Message:   "Invalid token",
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
			ctx := context.WithValue(context.Background(), models.JwtClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

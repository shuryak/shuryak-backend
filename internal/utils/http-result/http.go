package http_result

import (
	"encoding/json"
	"github.com/shuryak/shuryak-backend/internal/models"
	"net/http"
)

func WriteError(w *http.ResponseWriter, errorCode models.ErrorCode, description string) {
	var httpStatusCode int

	switch errorCode {
	case models.BadRequest:
		httpStatusCode = http.StatusBadRequest
	case models.InternalError:
		httpStatusCode = http.StatusInternalServerError
	case models.BadAuth:
		httpStatusCode = http.StatusBadRequest
	case models.NotAuthorized:
		httpStatusCode = http.StatusUnauthorized
	case models.InvalidToken:
		httpStatusCode = http.StatusUnauthorized
	case models.ExpiredToken:
		httpStatusCode = http.StatusForbidden
	case models.NotUniqueData:
		httpStatusCode = http.StatusBadRequest
	case models.InvalidFieldLength:
		httpStatusCode = http.StatusBadRequest
	default:
		httpStatusCode = http.StatusInternalServerError
	}

	(*w).WriteHeader(httpStatusCode)

	errorMessage := models.ErrorDTO{
		ErrorCode: errorCode,
		Message:   description,
	}
	json.NewEncoder(*w).Encode(errorMessage)
}

func WriteEmpty(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusOK)

	json.NewEncoder(*w).Encode(struct {
	}{})
}

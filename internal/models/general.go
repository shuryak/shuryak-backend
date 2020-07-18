package models

type ErrorCode int

const (
	BadRequest      ErrorCode = 0
	InternalError   ErrorCode = 1
	BadAuth         ErrorCode = 2
	ExpiredAuthData ErrorCode = 3
)

type ErrorDTO struct {
	ErrorCode ErrorCode `json:"error_code"`
	Message   string    `json:"message"`
}

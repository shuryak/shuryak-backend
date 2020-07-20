package models

type ErrorCode int

const (
	BadRequest         ErrorCode = 0
	InternalError      ErrorCode = 1
	BadAuth            ErrorCode = 2
	ExpiredAuthData    ErrorCode = 3
	NotUniqueData      ErrorCode = 4
	InvalidFieldLength ErrorCode = 5
)

type ErrorDTO struct {
	ErrorCode ErrorCode `json:"error_code"`
	Message   string    `json:"message"`
}

type FindExpression struct {
	Query string `json:"query"`
}

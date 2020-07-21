package models

type ErrorCode int

const (
	BadRequest         ErrorCode = 0 // Bad request (user error)
	InternalError      ErrorCode = 1 // Server error (╯°□°）╯︵ ┻━┻
	BadAuth            ErrorCode = 2 // Bad login details (user error)
	NotAuthorized      ErrorCode = 3 // To perform the action, you must pass an access token (user error)
	InvalidToken       ErrorCode = 4 // Invalid token (user error)
	ExpiredToken       ErrorCode = 5 // Expired token (user error)
	NotUniqueData      ErrorCode = 6 // Data is not unique when needed (user error)
	InvalidFieldLength ErrorCode = 7 // Invalid field length (user error)
)

type ErrorDTO struct {
	ErrorCode ErrorCode `json:"error_code"`
	Message   string    `json:"message"`
}

type FindExpression struct {
	Query string `json:"query"`
}

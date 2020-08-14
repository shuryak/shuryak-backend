package models

type ErrorCode int
type Limit int
type CtxKey uint

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

const (
	FirstNameMinLimit Limit = 2
	FirstNameMaxLimit Limit = 32
	LastNameMinLimit  Limit = 2
	LastNameMaxLimit  Limit = 32
	NicknameMinLimit  Limit = 2
	NicknameMaxLimit  Limit = 16
	PasswordMinLimit  Limit = 8
	PasswordMaxLimit  Limit = 256

	ArticleIdMinLimit   Limit = 3
	ArticleIdMaxLimit   Limit = 24
	ArticleNameMinLimit Limit = 3
	ArticleNameMaxLimit Limit = 100

	FindMaxLimit Limit = 10
)

const (
	JwtClaimsKey CtxKey = 0
)

type ErrorDTO struct {
	ErrorCode ErrorCode `json:"error_code"`
	Message   string    `json:"message"`
}

type FindOneExpression struct {
	Query string `json:"query"`
}

type FindManyExpression struct {
	Query  string `json:"query"`
	Count  uint   `json:"count"`
	Offset uint   `json:"offset"`
}

type GetListExpression struct {
	Count  uint `json:"count"`
	Offset uint `json:"offset"`
}

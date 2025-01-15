package errors

import "fmt"

type ErrorCode string

const (
	ErrInvalidToken   ErrorCode = "INVALID_TOKEN"
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrServerError    ErrorCode = "SERVER_ERROR"
	ErrNetworkError   ErrorCode = "NETWORK_ERROR"
)

type Error struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewError(code ErrorCode, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Package errors defines errors returned by the ZenMoney SDK.
package errors

import "fmt"

// ErrorCode identifies a category of SDK error.
type ErrorCode string

const (
	ErrInvalidToken     ErrorCode = "INVALID_TOKEN"
	ErrInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrServerError      ErrorCode = "SERVER_ERROR"
	ErrNetworkError     ErrorCode = "NETWORK_ERROR"
	ErrRateLimit        ErrorCode = "RATE_LIMIT"
	ErrResponseTooLarge ErrorCode = "RESPONSE_TOO_LARGE"
)

// Error describes an error returned by the SDK.
type Error struct {
	Code    ErrorCode
	Message string
	Err     error

	// StatusCode is the HTTP response status. It is zero for non-HTTP errors.
	StatusCode int
	// BodySnippet is a bounded fragment of an HTTP error response body.
	BodySnippet string
	// BodyTruncated reports whether BodySnippet omits bytes from the response body.
	BodyTruncated bool
	// RequestID identifies the failed request when the server provides an ID header.
	RequestID string
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error, if any.
func (e *Error) Unwrap() error {
	return e.Err
}

// New creates an SDK error with the provided category and cause.
func New(code ErrorCode, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

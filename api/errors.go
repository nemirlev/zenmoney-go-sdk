package api

import sdkerrors "github.com/nemirlev/zenmoney-go-sdk/v2/errors"

// ErrorCode identifies a category of SDK error.
type ErrorCode = sdkerrors.ErrorCode

const (
	ErrInvalidToken     = sdkerrors.ErrInvalidToken
	ErrInvalidRequest   = sdkerrors.ErrInvalidRequest
	ErrServerError      = sdkerrors.ErrServerError
	ErrNetworkError     = sdkerrors.ErrNetworkError
	ErrRateLimit        = sdkerrors.ErrRateLimit
	ErrResponseTooLarge = sdkerrors.ErrResponseTooLarge
)

// Error describes an error returned by the SDK.
type Error = sdkerrors.Error

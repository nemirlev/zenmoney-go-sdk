package api_test

import (
	stdErrors "errors"
	"testing"

	"github.com/nemirlev/zenmoney-go-sdk/v3/api"
	"github.com/stretchr/testify/require"
)

func TestNewClientReturnsPublicError(t *testing.T) {
	_, err := api.NewClient("")

	var apiErr *api.Error
	require.True(t, stdErrors.As(err, &apiErr))
	require.Equal(t, api.ErrInvalidToken, apiErr.Code)
}

package errors_test

import (
	stdErrors "errors"
	"testing"

	"github.com/nemirlev/zenmoney-go-sdk/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	t.Run("error without wrap", func(t *testing.T) {
		err := errors.NewError(errors.ErrInvalidToken, "token not provided", nil)
		require.NotNil(t, err)
		require.Equal(t, errors.ErrInvalidToken, err.Code)
		require.Equal(t, "INVALID_TOKEN: token not provided", err.Error())

		// Unwrap не должен ничего возвращать, т.к. мы не оборачивали внутреннюю ошибку
		require.Nil(t, stdErrors.Unwrap(err))
	})

	t.Run("error with wrap", func(t *testing.T) {
		rootErr := stdErrors.New("root error")
		err := errors.NewError(errors.ErrNetworkError, "network issue occurred", rootErr)
		require.NotNil(t, err)
		require.Equal(t, errors.ErrNetworkError, err.Code)
		require.Contains(t, err.Error(), "NETWORK_ERROR: network issue occurred: root error")

		// Проверим, что Unwrap() вернёт нашу изначальную ошибку
		unwrapped := stdErrors.Unwrap(err)
		require.Equal(t, rootErr, unwrapped)
	})
}

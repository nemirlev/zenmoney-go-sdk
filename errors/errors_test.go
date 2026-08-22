package errors_test

import (
	stdErrors "errors"
	"testing"

	sdkerrors "github.com/nemirlev/zenmoney-go-sdk/v3/errors"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := sdkerrors.New(sdkerrors.ErrInvalidToken, "token not provided", nil)

		require.Equal(t, sdkerrors.ErrInvalidToken, err.Code)
		require.Equal(t, "INVALID_TOKEN: token not provided", err.Error())
		require.Nil(t, stdErrors.Unwrap(err))
	})

	t.Run("with cause", func(t *testing.T) {
		cause := stdErrors.New("root error")
		err := sdkerrors.New(sdkerrors.ErrNetworkError, "network issue occurred", cause)

		require.ErrorIs(t, err, cause)
		require.Equal(t, "NETWORK_ERROR: network issue occurred: root error", err.Error())
	})
}

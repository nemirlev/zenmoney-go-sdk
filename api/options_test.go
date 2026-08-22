package api_test

import (
	"context"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClientOptionsValidation(t *testing.T) {
	var nilOption api.Option

	tests := []struct {
		name string
		opts []api.Option
	}{
		{name: "nil option", opts: []api.Option{nilOption}},
		{name: "nil HTTP client", opts: []api.Option{api.WithHTTPClient(nil)}},
		{name: "negative timeout", opts: []api.Option{api.WithTimeout(-time.Second)}},
		{name: "zero response limit", opts: []api.Option{api.WithMaxResponseSize(0)}},
		{name: "negative response limit", opts: []api.Option{api.WithMaxResponseSize(-1)}},
		{name: "overflowing response limit", opts: []api.Option{api.WithMaxResponseSize(math.MaxInt64)}},
		{name: "empty base URL", opts: []api.Option{api.WithBaseURL("")}},
		{name: "malformed base URL", opts: []api.Option{api.WithBaseURL("://invalid")}},
		{name: "relative base URL", opts: []api.Option{api.WithBaseURL("api.test.com/v8/")}},
		{name: "unsupported base URL scheme", opts: []api.Option{api.WithBaseURL("ftp://api.test.com/v8/")}},
		{name: "base URL query", opts: []api.Option{api.WithBaseURL("https://api.test.com/v8/?debug=true")}},
		{name: "base URL fragment", opts: []api.Option{api.WithBaseURL("https://api.test.com/v8/#diff")}},
		{
			name: "negative retry attempts",
			opts: []api.Option{api.WithRetryPolicy(-1, time.Second)},
		},
		{
			name: "negative retry wait",
			opts: []api.Option{api.WithRetryPolicy(1, -time.Second)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := api.NewClient("test-token", tt.opts...)

			require.Nil(t, client)
			var apiErr *api.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, api.ErrInvalidRequest, apiErr.Code)
		})
	}
}

func TestWithTimeoutLimitsRequest(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			<-req.Context().Done()
			return nil, req.Context().Err()
		}),
	}
	client, err := api.NewClient(
		"test-token",
		api.WithHTTPClient(httpClient),
		api.WithTimeout(20*time.Millisecond),
		api.WithRetryPolicy(3, time.Second),
	)
	require.NoError(t, err)
	started := time.Now()

	_, err = client.FullSync(context.Background())

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(started), 500*time.Millisecond)
}

func TestWithMaxResponseSizeLimitsResponse(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	client, err := api.NewClient(
		"test-token",
		api.WithHTTPClient(httpClient),
		api.WithMaxResponseSize(1),
	)
	require.NoError(t, err)

	_, err = client.FullSync(context.Background())

	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, api.ErrResponseTooLarge, apiErr.Code)
}

package api_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
	"github.com/stretchr/testify/require"
)

func TestForceSyncEntitiesSincePreservesCursor(t *testing.T) {
	lastSync := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var body models.Request
			require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
			require.Equal(t, int(lastSync.Unix()), body.ServerTimestamp)
			require.Equal(t, []models.EntityType{
				models.EntityTypeTransaction,
				models.EntityTypeAccount,
			}, body.ForceFetch)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(
					`{"serverTimestamp":1718454660}`,
				)),
				Header: make(http.Header),
			}, nil
		}),
	}
	client, err := api.NewClient(
		"test-token",
		api.WithHTTPClient(httpClient),
		api.WithRetryPolicy(0, 0),
	)
	require.NoError(t, err)

	response, err := client.ForceSyncEntitiesSince(
		context.Background(),
		lastSync,
		models.EntityTypeTransaction,
		models.EntityTypeAccount,
	)

	require.NoError(t, err)
	require.Equal(t, 1718454660, response.ServerTimestamp)
}

func TestForceSyncEntitiesCompatibilityWrapper(t *testing.T) {
	var requestBody models.Request
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			require.NoError(t, json.NewDecoder(req.Body).Decode(&requestBody))

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"serverTimestamp":1718454660}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}
	client, err := api.NewClient(
		"test-token",
		api.WithHTTPClient(httpClient),
		api.WithRetryPolicy(0, 0),
	)
	require.NoError(t, err)
	startedAt := time.Now()

	_, err = client.ForceSyncEntities( //nolint:staticcheck // Verify the v2 compatibility wrapper.
		context.Background(),
		models.EntityTypeTransaction,
	)

	require.NoError(t, err)
	require.Equal(t, []models.EntityType{models.EntityTypeTransaction}, requestBody.ForceFetch)
	require.WithinDuration(
		t,
		startedAt,
		time.Unix(int64(requestBody.ServerTimestamp), 0),
		time.Second,
	)
}

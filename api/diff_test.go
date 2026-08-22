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
			require.Equal(t, lastSync.Unix(), body.ServerTimestamp)
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
	require.Equal(t, int64(1718454660), response.ServerTimestamp)
}

func TestSyncSendsProvidedRequest(t *testing.T) {
	want := models.Request{
		CurrentClientTimestamp: 1_718_455_000,
		ServerTimestamp:        1_718_450_000,
		ForceFetch:             []models.EntityType{models.EntityTypeAccount},
	}
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, http.MethodPost, req.Method)
			require.Equal(t, "/v8/diff/", req.URL.Path)

			var got models.Request
			require.NoError(t, json.NewDecoder(req.Body).Decode(&got))
			require.Equal(t, want, got)

			return jsonResponse(`{"serverTimestamp":1718456000}`), nil
		}),
	}
	client, err := api.NewClient("test-token", api.WithHTTPClient(httpClient))
	require.NoError(t, err)

	response, err := client.Sync(context.Background(), want)

	require.NoError(t, err)
	require.Equal(t, int64(1_718_456_000), response.ServerTimestamp)
}

func TestSyncSinceUsesPreviousTimestamp(t *testing.T) {
	lastSync := time.Date(2025, 2, 3, 4, 5, 6, 0, time.UTC)
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var got models.Request
			require.NoError(t, json.NewDecoder(req.Body).Decode(&got))
			require.Equal(t, lastSync.Unix(), got.ServerTimestamp)
			require.Greater(t, got.CurrentClientTimestamp, int64(0))

			return jsonResponse(`{"serverTimestamp":1738555600}`), nil
		}),
	}
	client, err := api.NewClient("test-token", api.WithHTTPClient(httpClient))
	require.NoError(t, err)

	response, err := client.SyncSince(context.Background(), lastSync)

	require.NoError(t, err)
	require.Equal(t, int64(1_738_555_600), response.ServerTimestamp)
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
		time.Unix(requestBody.ServerTimestamp, 0),
		time.Second,
	)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

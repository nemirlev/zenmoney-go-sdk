package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/nemirlev/zenmoney-go-sdk/v3/api"
	"github.com/nemirlev/zenmoney-go-sdk/v3/models"
	"github.com/stretchr/testify/require"
)

func TestSuggestUsesPublicClient(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			require.Equal(t, http.MethodPost, req.Method)
			require.Equal(t, "/v8/suggest/", req.URL.Path)

			var got models.Transaction
			require.NoError(t, json.NewDecoder(req.Body).Decode(&got))
			require.Equal(t, "Coffee shop", got.Payee)

			return jsonResponse(`{"payee":"Coffee shop","tag":["food"]}`), nil
		}),
	}
	client, err := api.NewClient("test-token", api.WithHTTPClient(httpClient))
	require.NoError(t, err)

	suggestion, err := client.Suggest(context.Background(), models.Transaction{Payee: "Coffee shop"})

	require.NoError(t, err)
	require.Equal(t, "Coffee shop", suggestion.Payee)
	require.Equal(t, []string{"food"}, suggestion.Tag)
}

func TestSuggestBatchPreservesOrder(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var got []models.Transaction
			require.NoError(t, json.NewDecoder(req.Body).Decode(&got))
			require.Equal(t, []string{"First", "Second"}, []string{got[0].Payee, got[1].Payee})

			return jsonResponse(`[{"payee":"First","tag":["one"]},{"payee":"Second","tag":["two"]}]`), nil
		}),
	}
	client, err := api.NewClient("test-token", api.WithHTTPClient(httpClient))
	require.NoError(t, err)

	suggestions, err := client.SuggestBatch(context.Background(), []models.Transaction{
		{Payee: "First"},
		{Payee: "Second"},
	})

	require.NoError(t, err)
	require.Equal(t, []string{"First", "Second"}, []string{suggestions[0].Payee, suggestions[1].Payee})
}

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/internal/errors"
	"github.com/nemirlev/zenmoney-go-sdk/models"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return m.response, m.err
}

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)

	client, err := NewClient(
		"test-token",
		server.URL+"/",
		&http.Client{},
		time.Second,
		0, // Отключаем ретраи для тестов
		time.Second,
	)
	require.NoError(t, err)

	return server, client
}

func TestNewClient(t *testing.T) {
	t.Run("successfully creates client", func(t *testing.T) {
		client, err := NewClient(
			"test-token",
			"https://api.test.com/",
			&http.Client{},
			time.Second,
			3,
			time.Second,
		)
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("fails with empty token", func(t *testing.T) {
		client, err := NewClient(
			"",
			"https://api.test.com/",
			&http.Client{},
			time.Second,
			3,
			time.Second,
		)
		require.Error(t, err)
		require.Nil(t, client)
		require.Equal(t, errors.ErrInvalidToken, err.(*errors.Error).Code)
	})
}

func TestSync(t *testing.T) {
	t.Run("successful sync", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод и путь
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/diff/", r.URL.Path)

			// Проверяем заголовки
			require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Проверяем тело запроса
			var reqBody models.Request
			require.NoError(t, json.NewDecoder(r.Body).Decode(&reqBody))
			require.Greater(t, reqBody.CurrentClientTimestamp, 0)
			require.Equal(t, 1642300700, reqBody.ServerTimestamp)

			// Возвращаем ответ
			response := models.Response{
				ServerTimestamp: 1642300800,
				Instrument: []models.Instrument{
					{
						ID:         1,
						Title:      "US Dollar",
						ShortTitle: "USD",
						Symbol:     "$",
						Rate:       74.5,
						Changed:    1642300700,
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		})
		defer server.Close()

		resp, err := client.Sync(context.Background(), models.Request{
			CurrentClientTimestamp: int(time.Now().Unix()),
			ServerTimestamp:        1642300700,
		})

		require.NoError(t, err)
		require.Equal(t, 1642300800, resp.ServerTimestamp)
		require.Len(t, resp.Instrument, 1)
		require.Equal(t, "USD", resp.Instrument[0].ShortTitle)
	})

	t.Run("handles network error", func(t *testing.T) {
		client, err := NewClient(
			"test-token",
			"https://invalid.url/",
			&http.Client{},
			time.Second,
			0,
			time.Second,
		)
		require.NoError(t, err)

		_, err = client.Sync(context.Background(), models.Request{})
		require.Error(t, err)
		apiErr, ok := err.(*errors.Error)
		require.True(t, ok)
		require.Equal(t, errors.ErrNetworkError, apiErr.Code)
	})

	t.Run("handles server error", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		_, err := client.Sync(context.Background(), models.Request{})
		require.Error(t, err)
		apiErr, ok := err.(*errors.Error)
		require.True(t, ok)
		require.Equal(t, errors.ErrServerError, apiErr.Code)
	})

	t.Run("handles invalid JSON response", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`invalid json`))
		})
		defer server.Close()

		_, err := client.Sync(context.Background(), models.Request{})
		require.Error(t, err)
		apiErr, ok := err.(*errors.Error)
		require.True(t, ok)
		require.Equal(t, errors.ErrInvalidRequest, apiErr.Code)
	})
}

func TestSuggest(t *testing.T) {
	t.Run("successful suggestion", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/suggest/", r.URL.Path)

			// Проверяем тело запроса
			var reqTx models.Transaction
			require.NoError(t, json.NewDecoder(r.Body).Decode(&reqTx))
			require.Equal(t, "McDonalds", reqTx.Payee)

			// Возвращаем предложение
			suggestion := models.Transaction{
				Payee:    "McDonalds",
				Merchant: stringPtr("mcdonalds-1"),
				Tag:      []string{"food", "fast-food"},
			}
			json.NewEncoder(w).Encode(suggestion)
		})
		defer server.Close()

		tx := models.Transaction{
			Payee: "McDonalds",
		}

		suggestion, err := client.Suggest(context.Background(), tx)
		require.NoError(t, err)
		require.Equal(t, "McDonalds", suggestion.Payee)
		require.Equal(t, "mcdonalds-1", *suggestion.Merchant)
		require.Equal(t, []string{"food", "fast-food"}, suggestion.Tag)
	})

	t.Run("handles server error", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
		defer server.Close()

		_, err := client.Suggest(context.Background(), models.Transaction{})
		require.Error(t, err)
		apiErr, ok := err.(*errors.Error)
		require.True(t, ok)
		require.Equal(t, errors.ErrServerError, apiErr.Code)
	})
}

func TestSuggestBatch(t *testing.T) {
	t.Run("successful batch suggestion", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/suggest/", r.URL.Path)

			// Проверяем тело запроса
			var reqTxs []models.Transaction
			require.NoError(t, json.NewDecoder(r.Body).Decode(&reqTxs))
			require.Len(t, reqTxs, 2)
			require.Equal(t, "McDonalds", reqTxs[0].Payee)
			require.Equal(t, "Starbucks", reqTxs[1].Payee)

			// Возвращаем предложения
			suggestions := []models.Transaction{
				{
					Payee:    "McDonalds",
					Merchant: stringPtr("mcdonalds-1"),
					Tag:      []string{"food", "fast-food"},
				},
				{
					Payee:    "Starbucks",
					Merchant: stringPtr("starbucks-1"),
					Tag:      []string{"food", "coffee"},
				},
			}
			json.NewEncoder(w).Encode(suggestions)
		})
		defer server.Close()

		txs := []models.Transaction{
			{Payee: "McDonalds"},
			{Payee: "Starbucks"},
		}

		suggestions, err := client.SuggestBatch(context.Background(), txs)
		require.NoError(t, err)
		require.Len(t, suggestions, 2)
		require.Equal(t, "McDonalds", suggestions[0].Payee)
		require.Equal(t, "mcdonalds-1", *suggestions[0].Merchant)
		require.Equal(t, []string{"food", "fast-food"}, suggestions[0].Tag)
		require.Equal(t, "Starbucks", suggestions[1].Payee)
		require.Equal(t, "starbucks-1", *suggestions[1].Merchant)
		require.Equal(t, []string{"food", "coffee"}, suggestions[1].Tag)
	})

	t.Run("handles empty batch", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode([]models.Transaction{})
		})
		defer server.Close()

		suggestions, err := client.SuggestBatch(context.Background(), []models.Transaction{})
		require.NoError(t, err)
		require.Empty(t, suggestions)
	})
}
func TestFullSync(t *testing.T) {
	t.Run("successful full sync with correct request", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод и путь
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/diff/", r.URL.Path)

			// Проверяем body запроса
			var req models.Request
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)

			// В FullSync serverTimestamp должен быть 0
			require.Equal(t, 0, req.ServerTimestamp)
			// CurrentClientTimestamp должен быть текущим временем
			require.Greater(t, req.CurrentClientTimestamp, 0)
			require.LessOrEqual(t, req.CurrentClientTimestamp, int(time.Now().Unix()))

			// Возвращаем тестовый ответ
			resp := models.Response{
				ServerTimestamp: int(time.Now().Unix()),
				User: []models.User{
					{
						ID:      1,
						Login:   "testuser",
						Changed: int(time.Now().Unix()),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		resp, err := client.FullSync(context.Background())
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Greater(t, resp.ServerTimestamp, 0)
		require.Len(t, resp.User, 1)
		require.Equal(t, "testuser", resp.User[0].Login)
	})

	t.Run("handles server error", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		_, err := client.FullSync(context.Background())
		require.Error(t, err)
	})
}

func TestSyncSince(t *testing.T) {
	t.Run("successful sync since timestamp", func(t *testing.T) {
		lastSync := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/diff/", r.URL.Path)

			var req models.Request
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)

			// Проверяем, что timestamp соответствует переданному времени
			require.Equal(t, int(lastSync.Unix()), req.ServerTimestamp)
			require.Greater(t, req.CurrentClientTimestamp, int(lastSync.Unix()))

			resp := models.Response{
				ServerTimestamp: int(time.Now().Unix()),
				Transaction: []models.Transaction{
					{
						ID:      "test-tx-1",
						Changed: int(time.Now().Unix()),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		resp, err := client.SyncSince(context.Background(), lastSync)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Transaction, 1)
		require.Equal(t, "test-tx-1", resp.Transaction[0].ID)
	})

	t.Run("handles invalid response", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		})
		defer server.Close()

		_, err := client.SyncSince(context.Background(), time.Now())
		require.Error(t, err)
	})
}

func TestForceSyncEntities(t *testing.T) {
	t.Run("successful force sync with multiple entities", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/diff/", r.URL.Path)

			var req models.Request
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)

			// Проверяем запрошенные сущности
			require.Contains(t, req.ForceFetch, models.EntityTypeTransaction)
			require.Contains(t, req.ForceFetch, models.EntityTypeAccount)
			require.Len(t, req.ForceFetch, 2)

			resp := models.Response{
				ServerTimestamp: int(time.Now().Unix()),
				Transaction: []models.Transaction{
					{
						ID:      "forced-tx-1",
						Changed: int(time.Now().Unix()),
					},
				},
				Account: []models.Account{
					{
						ID:      "forced-acc-1",
						Changed: int(time.Now().Unix()),
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		resp, err := client.ForceSyncEntities(context.Background(),
			models.EntityTypeTransaction,
			models.EntityTypeAccount,
		)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Transaction, 1)
		require.Equal(t, "forced-tx-1", resp.Transaction[0].ID)
		require.Len(t, resp.Account, 1)
		require.Equal(t, "forced-acc-1", resp.Account[0].ID)
	})

	t.Run("successful force sync with empty entities", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			var req models.Request
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)
			require.Empty(t, req.ForceFetch)

			resp := models.Response{
				ServerTimestamp: int(time.Now().Unix()),
			}
			json.NewEncoder(w).Encode(resp)
		})
		defer server.Close()

		resp, err := client.ForceSyncEntities(context.Background())
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Greater(t, resp.ServerTimestamp, 0)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			// Имитируем долгий ответ
			time.Sleep(100 * time.Millisecond)
		})
		defer server.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := client.ForceSyncEntities(ctx, models.EntityTypeTransaction)
		require.Error(t, err)
	})
}

// Вспомогательная функция для создания указателя на string
func stringPtr(s string) *string {
	return &s
}

// Package client provides internal implementation of ZenMoney API client
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/errors"
	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
)

// Client represents internal implementation of ZenMoney API client
type Client struct {
	baseURL       string
	token         string
	httpClient    *http.Client
	timeout       time.Duration
	retryAttempts int
	retryWaitTime time.Duration
}

// NewClient creates a new instance of the internal API client
func NewClient(token string, baseURL string, httpClient *http.Client, timeout time.Duration, retryAttempts int, retryWaitTime time.Duration) (*Client, error) {
	if token == "" {
		return nil, errors.New(errors.ErrInvalidToken, "token is not provided", nil)
	}
	if httpClient == nil {
		return nil, errors.New(errors.ErrInvalidRequest, "HTTP client is nil", nil)
	}
	if timeout < 0 {
		return nil, errors.New(errors.ErrInvalidRequest, "timeout must not be negative", nil)
	}
	if retryAttempts < 0 {
		return nil, errors.New(errors.ErrInvalidRequest, "retry attempts must not be negative", nil)
	}
	if retryWaitTime < 0 {
		return nil, errors.New(errors.ErrInvalidRequest, "retry wait time must not be negative", nil)
	}

	return &Client{
		baseURL:       baseURL,
		token:         token,
		httpClient:    httpClient,
		timeout:       timeout,
		retryAttempts: retryAttempts,
		retryWaitTime: retryWaitTime,
	}, nil
}

// sendRequest sends an HTTP request to the specified endpoint with the given method and body
// It handles retries, timeouts, and response processing
func (c *Client) sendRequest(ctx context.Context, endpoint string, method string, body any) ([]byte, error) {
	if ctx == nil {
		return nil, errors.New(errors.ErrInvalidRequest, "context is nil", nil)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New(errors.ErrInvalidRequest, "failed to marshal request body", err)
	}

	requestCtx := ctx
	cancel := func() {}
	if c.timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, c.timeout)
	}
	defer cancel()

	for attempt := 0; attempt <= c.retryAttempts; attempt++ {
		req, err := http.NewRequestWithContext(
			requestCtx,
			method,
			c.baseURL+endpoint,
			bytes.NewReader(jsonBody),
		)
		if err != nil {
			return nil, errors.New(errors.ErrInvalidRequest, "failed to create request", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.token)

		resp, requestErr := c.httpClient.Do(req)
		if requestErr == nil {
			return readResponse(resp)
		}
		closeResponse(resp)

		if requestCtx.Err() != nil {
			return nil, errors.New(errors.ErrNetworkError, "request context ended", requestCtx.Err())
		}
		if attempt == c.retryAttempts {
			return nil, errors.New(errors.ErrNetworkError, "failed to send request after retries", requestErr)
		}

		if err := waitForRetry(requestCtx, c.retryWaitTime); err != nil {
			return nil, errors.New(errors.ErrNetworkError, "retry interrupted", err)
		}
	}

	panic("unreachable")
}

func waitForRetry(ctx context.Context, waitTime time.Duration) error {
	if waitTime == 0 {
		return ctx.Err()
	}

	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func readResponse(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New(errors.ErrNetworkError, "got nil response", nil)
	}
	if resp.Body == nil {
		return nil, errors.New(errors.ErrNetworkError, "got response with nil body", nil)
	}

	resBody, readErr := io.ReadAll(resp.Body)
	closeErr := resp.Body.Close()
	if readErr != nil {
		return nil, errors.New(errors.ErrNetworkError, "failed to read response body", readErr)
	}
	if closeErr != nil {
		return nil, errors.New(errors.ErrNetworkError, "failed to close response body", closeErr)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &errors.Error{
			Code:       errorCodeForStatus(resp.StatusCode),
			Message:    fmt.Sprintf("server returned error status: %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
		}
	}

	return resBody, nil
}

func closeResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
}

func errorCodeForStatus(statusCode int) errors.ErrorCode {
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.ErrInvalidToken
	case http.StatusTooManyRequests:
		return errors.ErrRateLimit
	default:
		if statusCode >= http.StatusInternalServerError {
			return errors.ErrServerError
		}
		return errors.ErrInvalidRequest
	}
}

// Sync sends a synchronization request to ZenMoney API with the provided parameters
func (c *Client) Sync(ctx context.Context, body models.Request) (models.Response, error) {
	resBody, err := c.sendRequest(ctx, "diff/", http.MethodPost, body)
	if err != nil {
		return models.Response{}, err
	}

	var result models.Response
	if err := json.Unmarshal(resBody, &result); err != nil {
		return models.Response{}, errors.New(errors.ErrInvalidRequest,
			"failed to unmarshal response", err)
	}

	return result, nil
}

// FullSync performs a full synchronization with ZenMoney API, retrieving all available data
func (c *Client) FullSync(ctx context.Context) (models.Response, error) {
	body := models.Request{
		CurrentClientTimestamp: int(time.Now().Unix()),
		ServerTimestamp:        0,
	}

	return c.Sync(ctx, body)
}

// SyncSince performs a synchronization with ZenMoney API from the specified timestamp
func (c *Client) SyncSince(ctx context.Context, lastSync time.Time) (models.Response, error) {
	body := models.Request{
		CurrentClientTimestamp: int(time.Now().Unix()),
		ServerTimestamp:        int(lastSync.Unix()),
	}

	return c.Sync(ctx, body)
}

// ForceSyncEntities requests a full update of the specified entities without an incremental cursor.
func (c *Client) ForceSyncEntities(ctx context.Context, entityTypes ...models.EntityType) (models.Response, error) {
	return c.ForceSyncEntitiesSince(ctx, time.Now(), entityTypes...)
}

// ForceSyncEntitiesSince requests all specified entities and regular changes
// since the server timestamp returned by the previous synchronization.
func (c *Client) ForceSyncEntitiesSince(ctx context.Context, lastSync time.Time, entityTypes ...models.EntityType) (models.Response, error) {
	body := models.Request{
		CurrentClientTimestamp: int(time.Now().Unix()),
		ServerTimestamp:        int(lastSync.Unix()),
		ForceFetch:             entityTypes,
	}

	return c.Sync(ctx, body)
}

// Suggest sends a suggestion request to the ZenMoney API for a single transaction.
// It sends a POST request to the suggest endpoint with the provided transaction data.
// Only the fields present in the input transaction will be considered for suggestions.
//
// Parameters
//   - ctx: Context for the request
//   - transaction: Transaction object, can be partially filled
//
// Returns:
//   - Transaction: Transaction object with suggested values
//   - error: Any error that occurred during the request
func (c *Client) Suggest(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	resBody, err := c.sendRequest(ctx, "suggest/", http.MethodPost, transaction)
	if err != nil {
		return models.Transaction{}, err
	}

	var result models.Transaction
	if err := json.Unmarshal(resBody, &result); err != nil {
		return models.Transaction{}, errors.New(errors.ErrInvalidRequest,
			"failed to unmarshal response", err)
	}

	return result, nil
}

// SuggestBatch sends a batch suggestion request to the ZenMoney API for multiple transactions.
func (c *Client) SuggestBatch(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, error) {
	resBody, err := c.sendRequest(ctx, "suggest/", http.MethodPost, transactions)
	if err != nil {
		return nil, err
	}

	var result []models.Transaction
	if err := json.Unmarshal(resBody, &result); err != nil {
		return nil, errors.New(errors.ErrInvalidRequest,
			"failed to unmarshal response", err)
	}

	return result, nil
}

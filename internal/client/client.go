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

	"github.com/nemirlev/zenmoney-go-sdk/internal/errors"
	"github.com/nemirlev/zenmoney-go-sdk/models"
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
		return nil, errors.NewError(errors.ErrInvalidToken, "token is not provided", nil)
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
func (c *Client) sendRequest(ctx context.Context, endpoint string, method string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.NewError(errors.ErrInvalidRequest, "failed to marshal request body", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.NewError(errors.ErrInvalidRequest, "failed to create request", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.token)

	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.retryAttempts; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err == nil {
			break
		}
		lastErr = err
		time.Sleep(c.retryWaitTime)
	}

	if lastErr != nil {
		return nil, errors.NewError(errors.ErrNetworkError, "failed to send request after retries", lastErr)
	}

	if resp == nil {
		return nil, errors.NewError(errors.ErrNetworkError, "got nil response", nil)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			if err == nil {
				err = errors.NewError(errors.ErrNetworkError, "failed to close response body", cerr)
			}
			fmt.Printf("failed to close response body: %v\n", cerr)
		}
	}()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewError(errors.ErrNetworkError, "failed to read response body", err)
	}

	if resp.StatusCode >= 400 {
		return nil, errors.NewError(errors.ErrServerError,
			fmt.Sprintf("server returned error status: %d", resp.StatusCode), nil)
	}

	return resBody, nil
}

// Sync sends a synchronization request to ZenMoney API with the provided parameters
func (c *Client) Sync(ctx context.Context, body models.Request) (models.Response, error) {
	resBody, err := c.sendRequest(ctx, "diff/", http.MethodPost, body)
	if err != nil {
		return models.Response{}, err
	}

	var result models.Response
	if err := json.Unmarshal(resBody, &result); err != nil {
		return models.Response{}, errors.NewError(errors.ErrInvalidRequest,
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

// Package client provides internal implementation of ZenMoney API client
package client

import (
	"bytes"
	"context"
	"encoding/json"
	stdErrors "errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/errors"
	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
)

const maxHTTPErrorBodySnippet int64 = 8 << 10

// Client represents internal implementation of ZenMoney API client
type Client struct {
	baseURL         *url.URL
	token           string
	httpClient      *http.Client
	timeout         time.Duration
	retryAttempts   int
	retryWaitTime   time.Duration
	maxResponseSize int64
	logger          *slog.Logger
}

// NewClient creates a new instance of the internal API client
func NewClient(token string, baseURL string, httpClient *http.Client, timeout time.Duration, retryAttempts int, retryWaitTime time.Duration, maxResponseSize int64, logger *slog.Logger) (*Client, error) {
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
	if maxResponseSize <= 0 || maxResponseSize == math.MaxInt64 {
		return nil, errors.New(errors.ErrInvalidRequest, "maximum response size must be positive and less than math.MaxInt64", nil)
	}
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.New(errors.ErrInvalidRequest, "base URL is invalid", err)
	}
	if parsedBaseURL.Scheme != "http" && parsedBaseURL.Scheme != "https" {
		return nil, errors.New(errors.ErrInvalidRequest, "base URL must use HTTP or HTTPS", nil)
	}
	if parsedBaseURL.Host == "" {
		return nil, errors.New(errors.ErrInvalidRequest, "base URL must be absolute and include a host", nil)
	}
	if parsedBaseURL.RawQuery != "" || parsedBaseURL.Fragment != "" {
		return nil, errors.New(errors.ErrInvalidRequest, "base URL must not include a query or fragment", nil)
	}
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	return &Client{
		baseURL:         parsedBaseURL,
		token:           token,
		httpClient:      httpClient,
		timeout:         timeout,
		retryAttempts:   retryAttempts,
		retryWaitTime:   retryWaitTime,
		maxResponseSize: maxResponseSize,
		logger:          logger,
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
		attemptNumber := attempt + 1
		maxAttempts := c.retryAttempts + 1
		requestURL := c.baseURL.JoinPath(endpoint)
		req, err := http.NewRequestWithContext(
			requestCtx,
			method,
			requestURL.String(),
			bytes.NewReader(jsonBody),
		)
		if err != nil {
			return nil, errors.New(errors.ErrInvalidRequest, "failed to create request", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.token)

		startedAt := time.Now()
		c.logger.LogAttrs(
			requestCtx,
			slog.LevelDebug,
			"ZenMoney HTTP request started",
			slog.String("method", method),
			slog.String("endpoint", endpoint),
			slog.Int("attempt", attemptNumber),
			slog.Int("max_attempts", maxAttempts),
		)

		resp, requestErr := c.httpClient.Do(req)
		if requestErr == nil {
			resBody, responseErr := readResponse(resp, c.maxResponseSize)
			attrs := []slog.Attr{
				slog.String("method", method),
				slog.String("endpoint", endpoint),
				slog.Int("attempt", attemptNumber),
				slog.Int("max_attempts", maxAttempts),
				slog.Duration("duration", time.Since(startedAt)),
				slog.Int("status_code", responseStatusCode(resp)),
				slog.String("request_id", responseRequestID(responseHeader(resp))),
				slog.String("outcome", "success"),
			}
			if responseErr != nil {
				attrs[len(attrs)-1] = slog.String("outcome", "error")
				var sdkErr *errors.Error
				if stdErrors.As(responseErr, &sdkErr) {
					attrs = append(attrs, slog.String("error_code", string(sdkErr.Code)))
				}
			}
			c.logger.LogAttrs(requestCtx, slog.LevelDebug, "ZenMoney HTTP request completed", attrs...)

			return resBody, responseErr
		}
		closeResponse(resp)

		if requestCtx.Err() != nil {
			c.logTransportFailure(requestCtx, method, endpoint, attemptNumber, maxAttempts, startedAt, "context_ended")
			return nil, errors.New(errors.ErrNetworkError, "request context ended", requestCtx.Err())
		}
		if attempt == c.retryAttempts {
			c.logTransportFailure(requestCtx, method, endpoint, attemptNumber, maxAttempts, startedAt, "error")
			return nil, errors.New(errors.ErrNetworkError, "failed to send request after retries", requestErr)
		}

		c.logger.LogAttrs(
			requestCtx,
			slog.LevelWarn,
			"ZenMoney HTTP request retry scheduled",
			slog.String("method", method),
			slog.String("endpoint", endpoint),
			slog.Int("attempt", attemptNumber),
			slog.Int("max_attempts", maxAttempts),
			slog.Duration("retry_wait", c.retryWaitTime),
		)

		if err := waitForRetry(requestCtx, c.retryWaitTime); err != nil {
			return nil, errors.New(errors.ErrNetworkError, "retry interrupted", err)
		}
	}

	panic("unreachable")
}

func (c *Client) logTransportFailure(ctx context.Context, method string, endpoint string, attempt int, maxAttempts int, startedAt time.Time, outcome string) {
	c.logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		"ZenMoney HTTP request completed",
		slog.String("method", method),
		slog.String("endpoint", endpoint),
		slog.Int("attempt", attempt),
		slog.Int("max_attempts", maxAttempts),
		slog.Duration("duration", time.Since(startedAt)),
		slog.String("outcome", outcome),
		slog.String("error_code", string(errors.ErrNetworkError)),
	)
}

func responseStatusCode(resp *http.Response) int {
	if resp == nil {
		return 0
	}

	return resp.StatusCode
}

func responseHeader(resp *http.Response) http.Header {
	if resp == nil {
		return nil
	}

	return resp.Header
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

func readResponse(resp *http.Response, maxResponseSize int64) ([]byte, error) {
	if resp == nil {
		return nil, errors.New(errors.ErrNetworkError, "got nil response", nil)
	}
	if resp.Body == nil {
		return nil, errors.New(errors.ErrNetworkError, "got response with nil body", nil)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, readHTTPError(resp)
	}

	resBody, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
	closeErr := resp.Body.Close()
	if readErr != nil {
		return nil, errors.New(errors.ErrNetworkError, "failed to read response body", readErr)
	}
	if int64(len(resBody)) > maxResponseSize {
		return nil, errors.New(errors.ErrResponseTooLarge, "response body exceeds configured limit", nil)
	}
	if closeErr != nil {
		return nil, errors.New(errors.ErrNetworkError, "failed to close response body", closeErr)
	}

	return resBody, nil
}

func readHTTPError(resp *http.Response) error {
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxHTTPErrorBodySnippet+1))
	truncated := int64(len(body)) > maxHTTPErrorBodySnippet
	if truncated {
		body = body[:maxHTTPErrorBodySnippet]
	}
	closeErr := resp.Body.Close()

	cause := readErr
	if cause == nil {
		cause = closeErr
	}

	return &errors.Error{
		Code:          errorCodeForStatus(resp.StatusCode),
		Message:       fmt.Sprintf("server returned error status: %d", resp.StatusCode),
		Err:           cause,
		StatusCode:    resp.StatusCode,
		BodySnippet:   strings.ToValidUTF8(string(body), "\uFFFD"),
		BodyTruncated: truncated,
		RequestID:     responseRequestID(resp.Header),
	}
}

func responseRequestID(header http.Header) string {
	for _, name := range []string{"X-Request-ID", "Request-ID"} {
		if value := header.Get(name); value != "" {
			return value
		}
	}

	return ""
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
		CurrentClientTimestamp: time.Now().Unix(),
		ServerTimestamp:        0,
	}

	return c.Sync(ctx, body)
}

// SyncSince performs a synchronization with ZenMoney API from the specified timestamp
func (c *Client) SyncSince(ctx context.Context, lastSync time.Time) (models.Response, error) {
	body := models.Request{
		CurrentClientTimestamp: time.Now().Unix(),
		ServerTimestamp:        lastSync.Unix(),
	}

	return c.Sync(ctx, body)
}

// ForceSyncEntitiesSince requests all specified entities and regular changes
// since the server timestamp returned by the previous synchronization.
func (c *Client) ForceSyncEntitiesSince(ctx context.Context, lastSync time.Time, entityTypes ...models.EntityType) (models.Response, error) {
	body := models.Request{
		CurrentClientTimestamp: time.Now().Unix(),
		ServerTimestamp:        lastSync.Unix(),
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

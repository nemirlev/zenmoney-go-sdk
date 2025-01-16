package api

import (
	"context"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/internal/client"
	"github.com/nemirlev/zenmoney-go-sdk/models"
)

type Client struct {
	internal *client.Client
}

// NewClient creates a new instance of the ZenMoney API client
// token: authentication token for the API
// opts: optional configuration settings
// Returns: a new Client instance and an error if any
func NewClient(token string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	internalClient, err := client.NewClient(
		token,
		cfg.baseURL,
		cfg.httpClient,
		cfg.timeout,
		cfg.retryAttempts,
		cfg.retryWaitTime,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		internal: internalClient,
	}, nil
}

// Sync sends a synchronization request to ZenMoney API with the provided parameters
// ctx: context for the request
// body: request body containing synchronization parameters
// Returns: response from the API and an error if any
func (c *Client) Sync(ctx context.Context, body models.Request) (models.Response, error) {
	return c.internal.Sync(ctx, body)
}

// FullSync performs a full synchronization with ZenMoney API, retrieving all available data
// ctx: context for the request
// Returns: response from the API and an error if any
func (c *Client) FullSync(ctx context.Context) (models.Response, error) {
	return c.internal.FullSync(ctx)
}

// SyncSince performs a synchronization with ZenMoney API from the specified timestamp
// ctx: context for the request
// lastSync: timestamp from which to start synchronization
// Returns: response from the API and an error if any
func (c *Client) SyncSince(ctx context.Context, lastSync time.Time) (models.Response, error) {
	return c.internal.SyncSince(ctx, lastSync)
}

// ForceSyncEntities requests a full update of the specified entities along with regular changes
func (c *Client) ForceSyncEntities(ctx context.Context, entityTypes ...models.EntityType) (models.Response, error) {
	return c.internal.ForceSyncEntities(ctx, entityTypes...)
}

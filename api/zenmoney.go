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

func (c *Client) Sync(ctx context.Context, body models.Request) (models.Response, error) {
	return c.internal.Sync(ctx, body)
}

func (c *Client) FullSync(ctx context.Context) (models.Response, error) {
	return c.internal.FullSync(ctx)
}

func (c *Client) SyncSince(ctx context.Context, lastSync time.Time) (models.Response, error) {
	return c.internal.SyncSince(ctx, lastSync)
}

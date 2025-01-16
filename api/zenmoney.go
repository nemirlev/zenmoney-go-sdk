package api

import (
	"github.com/nemirlev/zenmoney-go-sdk/internal/client"
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

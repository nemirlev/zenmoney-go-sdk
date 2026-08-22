package api

import (
	"github.com/nemirlev/zenmoney-go-sdk/v3/errors"
	"github.com/nemirlev/zenmoney-go-sdk/v3/internal/client"
)

// Client provides access to the ZenMoney synchronization and suggestion APIs.
// A Client is safe for concurrent use. Its methods delegate transport details to
// an internal implementation so the public API can remain stable.
type Client struct {
	internal *client.Client
}

// NewClient creates a ZenMoney API client authenticated with token. Options are
// applied in order, so later options override earlier ones. Without options, the
// client uses the production API endpoint, a 30-second operation timeout, three
// transport retries, a 64 MiB response limit, and disabled diagnostics.
//
// NewClient returns an *Error when the token or any resulting configuration is
// invalid.
func NewClient(token string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New(errors.ErrInvalidRequest, "nil client option", nil)
		}
		opt(cfg)
	}

	internalClient, err := client.NewClient(
		token,
		cfg.baseURL,
		cfg.httpClient,
		cfg.timeout,
		cfg.retryAttempts,
		cfg.retryWaitTime,
		cfg.maxResponseSize,
		cfg.logger,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		internal: internalClient,
	}, nil
}

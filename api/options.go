package api

import (
	"log/slog"
	"net/http"
	"time"
)

// DefaultMaxResponseSize is the default maximum successful response body size.
const DefaultMaxResponseSize int64 = 64 << 20

// Config holds client configuration settings
type Config struct {
	baseURL         string
	httpClient      *http.Client
	timeout         time.Duration
	retryAttempts   int
	retryWaitTime   time.Duration
	maxResponseSize int64
	logger          *slog.Logger
}

// Option represents a function for configuring the client
type Option func(*Config)

// defaultConfig returns default configuration settings
func defaultConfig() *Config {
	return &Config{
		baseURL:         "https://api.zenmoney.ru/v8/",
		httpClient:      &http.Client{},
		timeout:         30 * time.Second,
		retryAttempts:   3,
		retryWaitTime:   1 * time.Second,
		maxResponseSize: DefaultMaxResponseSize,
		logger:          slog.New(slog.DiscardHandler),
	}
}

// WithMaxResponseSize limits successful response bodies to maxBytes.
// maxBytes must be positive and less than the maximum int64 value.
func WithMaxResponseSize(maxBytes int64) Option {
	return func(c *Config) {
		c.maxResponseSize = maxBytes
	}
}

// WithBaseURL sets the base URL used for API requests. The URL must be absolute,
// use HTTP or HTTPS, and must not contain a query or fragment.
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.baseURL = url
	}
}

// WithHTTPClient sets the HTTP client used to send requests. The client must not
// be nil. A custom Transport can be used to add tracing, metrics, or other HTTP
// middleware.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.httpClient = client
	}
}

// WithTimeout limits a complete SDK operation, including all retry attempts and
// waits between them. A zero timeout disables the SDK-level deadline; a negative
// timeout is rejected by NewClient.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WithRetryPolicy configures retries for transport failures. attempts is the
// number of retries after the initial request, and waitTime is the delay between
// attempts. HTTP error responses are not retried automatically.
func WithRetryPolicy(attempts int, waitTime time.Duration) Option {
	return func(c *Config) {
		c.retryAttempts = attempts
		c.retryWaitTime = waitTime
	}
}

// WithLogger enables optional structured diagnostics using logger. The SDK logs
// request metadata, outcomes, and retry decisions, but never authorization
// headers or request and response bodies. Passing nil disables diagnostics.
// Logging levels and output formatting are configured by the logger's handler.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}

package api

import (
	"net/http"
	"time"
)

// Config holds client configuration settings
type Config struct {
	baseURL       string
	httpClient    *http.Client
	timeout       time.Duration
	retryAttempts int
	retryWaitTime time.Duration
}

// Option represents a function for configuring the client
type Option func(*Config)

// defaultConfig returns default configuration settings
func defaultConfig() *Config {
	return &Config{
		baseURL:       "https://api.zenmoney.ru/v8/",
		httpClient:    &http.Client{},
		timeout:       30 * time.Second,
		retryAttempts: 3,
		retryWaitTime: 1 * time.Second,
	}
}

// WithBaseURL sets the base URL for API requests
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.httpClient = client
	}
}

// WithTimeout sets the timeout for API requests
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WithRetryPolicy sets the retry policy for failed requests
// attempts: number of retry attempts
// waitTime: duration to wait between retries
func WithRetryPolicy(attempts int, waitTime time.Duration) Option {
	return func(c *Config) {
		c.retryAttempts = attempts
		c.retryWaitTime = waitTime
	}
}

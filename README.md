# ZenMoney API GO SDK.

[![Go Reference](https://pkg.go.dev/badge/github.com/nemirlev/zenmoney-go-sdk/v2.svg)](https://pkg.go.dev/github.com/nemirlev/zenmoney-go-sdk/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/nemirlev/zenmoney-go-sdk/v2)](https://goreportcard.com/report/github.com/nemirlev/zenmoney-go-sdk/v2)
![GitHub License](https://img.shields.io/github/license/nemirlev/zenmoney-go-sdk)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/nemirlev/zenmoney-go-sdk)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/nemirlev/zenmoney-go-sdk)
[![codecov](https://codecov.io/gh/nemirlev/zenmoney-go-sdk/graph/badge.svg?token=J2S3N967Q7)](https://codecov.io/gh/nemirlev/zenmoney-go-sdk)

A robust and easy-to-use Go SDK for interacting with the ZenMoney API. This SDK provides a type-safe way to work with
ZenMoney's financial data synchronization API, including accounts, transactions, budgets, and more.

## Features

- 🚀 Easy-to-use, idiomatic Go API
- 🔒 Built-in retry mechanism with configurable policies
- 🔎 Optional structured diagnostics with `log/slog`
- 💪 Full type safety for all ZenMoney entities
- 🛡️ Comprehensive error handling
- ⚡ Support for all ZenMoney API operations
- 📚 Extensive documentation and examples

## Installation

```bash
go get github.com/nemirlev/zenmoney-go-sdk/v2
```

## Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
)

func main() {
	// Create a new client
	client, err := api.NewClient("your-token-here")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Perform full sync
	ctx := context.Background()
	resp, err := client.FullSync(ctx)
	if err != nil {
		log.Fatalf("Sync failed: %v", err)
	}

	// Work with the data
	for _, account := range resp.Account {
		if account.Balance != nil {
			log.Printf("Account: %s, Balance: %.2f", account.Title, *account.Balance)
		}
	}
}
```

## Configuration Options

The SDK supports various configuration options through the functional options pattern:

```go
client, err := api.NewClient(
    "your-token-here",
    api.WithTimeout(45*time.Second),
    api.WithRetryPolicy(5, 2*time.Second),
    api.WithMaxResponseSize(128 << 20), // 128 MiB
    api.WithBaseURL("https://custom-api.zenmoney.ru/v8/"),
)
```

`WithTimeout` limits the complete SDK operation, including retry waits. Retry
attempts are made only for transport errors; HTTP responses are returned as
typed API errors.

Successful response bodies are limited to 64 MiB by default. Use
`WithMaxResponseSize` when a full synchronization is expected to be larger.

### Structured diagnostics

Diagnostics are disabled by default. Pass a standard `slog.Logger` to inspect
request timing, status codes, request IDs, and transport retries:

```go
logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client, err := api.NewClient(
    "your-token-here",
    api.WithLogger(logger),
)
```

The logger's handler controls levels, formatting, and output. The SDK never logs
authorization headers or request and response bodies. `DEBUG` records describe
request starts and outcomes; `WARN` records describe scheduled retries.

## Available Operations

- Full synchronization
- Incremental synchronization from a specific timestamp
- Force sync specific entities from a known server timestamp
- Custom sync with specific parameters
- Suggestions for categories and operation merchants

## Error Handling

The SDK provides structured error types for better error handling:

```go
if err != nil {
    var apiErr *api.Error
    if errors.As(err, &apiErr) {
        switch apiErr.Code {
        case api.ErrInvalidToken:
            // Handle authentication error
        case api.ErrRateLimit:
            // Respect the API rate limit
        case api.ErrServerError:
            // Handle server error
        case api.ErrNetworkError:
            // Handle network error
        case api.ErrResponseTooLarge:
            // Raise the configured response limit if the payload is expected
        }
    }
}
```

HTTP errors also expose `StatusCode`, `RequestID`, `BodySnippet`, and
`BodyTruncated`. The response body fragment is limited to 8 KiB and is not
included in the error string.

## Examples

Check out the [examples](./examples) directory for more detailed usage examples:

- Basic usage
- Configuration options
- Sync operations
- Error handling
- Working with budgets
- Working with merchants

## API Documentation

For detailed API documentation, visit
the [Go package documentation](https://pkg.go.dev/github.com/nemirlev/zenmoney-go-sdk/v2).

For ZenMoney API documentation, visit
the [official API documentation](https://github.com/zenmoney/ZenPlugins/wiki/ZenMoney-API).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the ZenMoney team for providing the API
- Inspired by other excellent Go SDKs in the community

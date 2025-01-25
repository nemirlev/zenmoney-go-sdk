package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
	zerrors "github.com/nemirlev/zenmoney-go-sdk/v2/internal/errors"
)

func main() {
	client, err := api.NewClient("your-token-here")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example of handling different error types
	_, err = client.FullSync(ctx)
	if err != nil {
		var apiErr *zerrors.Error
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case zerrors.ErrInvalidToken:
				log.Printf("Authentication failed: %v", err)
				// Handle token refresh

			case zerrors.ErrServerError:
				log.Printf("Server error: %v", err)
				// Implement retry with backoff
				retryWithBackoff(client, ctx)

			case zerrors.ErrNetworkError:
				log.Printf("Network error: %v", err)
				// Check connectivity

			default:
				log.Printf("Unknown error: %v", err)
			}
		}
	}
}

func retryWithBackoff(client *api.Client, ctx context.Context) {
	backoff := time.Second
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		_, err := client.FullSync(ctx)
		if err == nil {
			return
		}

		log.Printf("Retry %d failed: %v", i+1, err)
		time.Sleep(backoff)
		backoff *= 2
	}
}

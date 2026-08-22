package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v3/api"
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
		if apiErr, ok := errors.AsType[*api.Error](err); ok {
			switch apiErr.Code {
			case api.ErrInvalidToken:
				log.Printf("Authentication failed: %v", err)
				// Handle token refresh

			case api.ErrServerError:
				log.Printf("Server error: %v", err)
				// Implement retry with backoff
				retryWithBackoff(client, ctx)

			case api.ErrNetworkError:
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

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/api"
)

func main() {
	// Example 1: Client with custom timeout and retry policy
	client1, err := api.NewClient(
		"your-token-here",
		api.WithTimeout(45*time.Second),
		api.WithRetryPolicy(5, 2*time.Second),
	)
	if err != nil {
		log.Printf("Failed to create client1: %v", err)
	}

	// Example 2: Client with custom HTTP client and base URL
	client2, err := api.NewClient(
		"your-token-here",
		api.WithBaseURL("https://custom-api.zenmoney.ru/v8/"),
		api.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
	)
	if err != nil {
		log.Printf("Failed to create client2: %v", err)
	}

	// Use clients...
	_ = client1
	_ = client2

	println("Clients created successfully")
}

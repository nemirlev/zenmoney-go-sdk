package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
)

func main() {
	client, err := api.NewClient("your-token-here")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Sync merchants
	resp, err := client.ForceSyncEntities(ctx, models.EntityTypeMerchant)
	if err != nil {
		log.Fatalf("Failed to sync merchants: %v", err)
	}

	// Process and display merchants
	for _, merchant := range resp.Merchant {
		fmt.Printf("Merchant: %s\n", merchant.Title)
		fmt.Printf("  ID: %s\n", merchant.ID)
		fmt.Printf("  Last changed: %v\n",
			time.Unix(int64(merchant.Changed), 0).Format(time.RFC3339))
	}
}

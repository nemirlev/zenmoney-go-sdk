package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nemirlev/zenmoney-go-sdk/v3/api"
)

func main() {
	// Create client with default configuration
	client, err := api.NewClient("your-token-here")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Perform full synchronization
	ctx := context.Background()
	resp, err := client.FullSync(ctx)
	if err != nil {
		log.Fatalf("Sync failed: %v", err)
	}

	// Print accounts
	for _, account := range resp.Account {
		if account.Balance != nil {
			fmt.Printf("Account: %s, Balance: %.2f\n", account.Title, *account.Balance)
		}
	}

	// Print recent transactions
	for _, tx := range resp.Transaction {
		fmt.Printf("Transaction: Date: %s, Amount: %.2f\n", tx.Date, tx.Income-tx.Outcome)
	}
}

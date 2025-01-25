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

	// Get current month's budgets
	resp, err := client.ForceSyncEntities(ctx, models.EntityTypeBudget)
	if err != nil {
		log.Fatalf("Failed to sync budgets: %v", err)
	}

	// Process budgets
	currentMonth := time.Now().Format("2006-01")
	for _, budget := range resp.Budget {
		if budget.Date[:7] == currentMonth {
			fmt.Printf("Budget for %s:\n", budget.Date)
			fmt.Printf("  Income: %.2f (Locked: %v)\n", budget.Income, budget.IncomeLock)
			fmt.Printf("  Outcome: %.2f (Locked: %v)\n", budget.Outcome, budget.OutcomeLock)

			if budget.Tag != nil {
				fmt.Printf("  Category key: %s\n", *budget.Tag)
			}
		}
	}
}

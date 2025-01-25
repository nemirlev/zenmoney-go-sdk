package main

import (
	"context"
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

	// Example 1: Full sync
	fullSyncResp, err := client.FullSync(ctx)
	if err != nil {
		log.Printf("Full sync failed: %v", err)
	}

	log.Printf("Full sync completed successfully")

	// Example 2: Sync since timestamp
	lastSync := time.Unix(int64(fullSyncResp.ServerTimestamp), 0)
	syncResp, err := client.SyncSince(ctx, lastSync)
	if err != nil {
		log.Printf("Sync since failed: %v", err)
	}

	log.Printf("Sync since completed successfully")

	// Example 3: Force sync specific entities
	forceSyncResp, err := client.ForceSyncEntities(ctx,
		models.EntityTypeAccount,
		models.EntityTypeTransaction,
	)
	if err != nil {
		log.Printf("Force sync failed: %v", err)
	}

	log.Printf("Force sync completed successfully")

	// Process responses...
	_ = fullSyncResp
	_ = syncResp
	_ = forceSyncResp

	log.Printf("All operations completed successfully")
}

package api_test

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v2/api"
	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
)

func ExampleNewClient() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	client, err := api.NewClient(
		os.Getenv("ZENMONEY_TOKEN"),
		api.WithTimeout(45*time.Second),
		api.WithRetryPolicy(2, time.Second),
		api.WithLogger(logger),
	)
	if err != nil {
		log.Print(err)
		return
	}

	_ = client
}

func ExampleClient_FullSync() {
	client, err := api.NewClient(os.Getenv("ZENMONEY_TOKEN"))
	if err != nil {
		log.Print(err)
		return
	}

	response, err := client.FullSync(context.Background())
	if err != nil {
		log.Print(err)
		return
	}

	_ = response.ServerTimestamp
}

func ExampleClient_Sync() {
	client, err := api.NewClient(os.Getenv("ZENMONEY_TOKEN"))
	if err != nil {
		log.Print(err)
		return
	}

	response, err := client.Sync(context.Background(), models.Request{
		CurrentClientTimestamp: time.Now().Unix(),
		ServerTimestamp:        0,
	})
	if err != nil {
		log.Print(err)
		return
	}

	_ = response.ServerTimestamp
}

func ExampleClient_SyncSince() {
	client, err := api.NewClient(os.Getenv("ZENMONEY_TOKEN"))
	if err != nil {
		log.Print(err)
		return
	}

	lastSync := time.Unix(1_700_000_000, 0)
	response, err := client.SyncSince(context.Background(), lastSync)
	if err != nil {
		log.Print(err)
		return
	}

	_ = response.ServerTimestamp
}

func ExampleClient_ForceSyncEntitiesSince() {
	client, err := api.NewClient(os.Getenv("ZENMONEY_TOKEN"))
	if err != nil {
		log.Print(err)
		return
	}

	lastSync := time.Unix(1_700_000_000, 0)
	response, err := client.ForceSyncEntitiesSince(
		context.Background(),
		lastSync,
		models.EntityTypeAccount,
		models.EntityTypeTransaction,
	)
	if err != nil {
		log.Print(err)
		return
	}

	_ = response.ServerTimestamp
}

func ExampleClient_Suggest() {
	client, err := api.NewClient(os.Getenv("ZENMONEY_TOKEN"))
	if err != nil {
		log.Print(err)
		return
	}

	suggestion, err := client.Suggest(context.Background(), models.Transaction{
		Payee: "Coffee shop",
	})
	if err != nil {
		log.Print(err)
		return
	}

	_ = suggestion
}

func ExampleClient_SuggestBatch() {
	client, err := api.NewClient(os.Getenv("ZENMONEY_TOKEN"))
	if err != nil {
		log.Print(err)
		return
	}

	suggestions, err := client.SuggestBatch(context.Background(), []models.Transaction{
		{Payee: "Coffee shop"},
		{Payee: "Grocery store"},
	})
	if err != nil {
		log.Print(err)
		return
	}

	_ = suggestions
}

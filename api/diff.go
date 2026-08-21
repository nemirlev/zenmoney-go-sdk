package api

import (
	"context"
	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
	"time"
)

// Sync sends a synchronization request to ZenMoney API with the provided parameters
// ctx: context for the request
// body: request body containing synchronization parameters
// Returns: response from the API and an error if any
func (c *Client) Sync(ctx context.Context, body models.Request) (models.Response, error) {
	return c.internal.Sync(ctx, body)
}

// FullSync performs a full synchronization with ZenMoney API, retrieving all available data
// ctx: context for the request
// Returns: response from the API and an error if any
func (c *Client) FullSync(ctx context.Context) (models.Response, error) {
	return c.internal.FullSync(ctx)
}

// SyncSince performs a synchronization with ZenMoney API from the specified timestamp
// ctx: context for the request
// lastSync: timestamp from which to start synchronization
// Returns: response from the API and an error if any
func (c *Client) SyncSince(ctx context.Context, lastSync time.Time) (models.Response, error) {
	return c.internal.SyncSince(ctx, lastSync)
}

// ForceSyncEntities requests a full update of the specified entities without an incremental cursor.
//
// Deprecated: Use ForceSyncEntitiesSince to include changes since a known server timestamp.
func (c *Client) ForceSyncEntities(ctx context.Context, entityTypes ...models.EntityType) (models.Response, error) {
	return c.internal.ForceSyncEntitiesSince(ctx, time.Now(), entityTypes...)
}

// ForceSyncEntitiesSince requests all specified entities and regular changes
// since the server timestamp returned by the previous synchronization.
func (c *Client) ForceSyncEntitiesSince(ctx context.Context, lastSync time.Time, entityTypes ...models.EntityType) (models.Response, error) {
	return c.internal.ForceSyncEntitiesSince(ctx, lastSync, entityTypes...)
}

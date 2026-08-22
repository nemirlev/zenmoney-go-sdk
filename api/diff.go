package api

import (
	"context"
	"time"

	"github.com/nemirlev/zenmoney-go-sdk/v3/models"
)

// Sync sends body to the ZenMoney diff endpoint and returns the decoded
// synchronization response. Use this method when the convenience methods do not
// provide enough control over timestamps, force-fetch entities, or changed data.
// The request uses ctx for cancellation and deadlines.
func (c *Client) Sync(ctx context.Context, body models.Request) (models.Response, error) {
	return c.internal.Sync(ctx, body)
}

// FullSync retrieves all data available to the authenticated user. It sends a
// zero server timestamp and uses the current time as the client timestamp. Large
// accounts may need a higher response limit configured with WithMaxResponseSize.
func (c *Client) FullSync(ctx context.Context) (models.Response, error) {
	return c.internal.FullSync(ctx)
}

// SyncSince retrieves changes since lastSync. Callers should persist the server
// timestamp returned by a previous response and pass it back as lastSync. The
// time is converted to Unix seconds before it is sent.
func (c *Client) SyncSince(ctx context.Context, lastSync time.Time) (models.Response, error) {
	return c.internal.SyncSince(ctx, lastSync)
}

// ForceSyncEntities requests a full update of entityTypes without accepting an
// explicit incremental cursor. It is retained for source compatibility.
//
// Deprecated: Use ForceSyncEntitiesSince to include changes since a known server timestamp.
func (c *Client) ForceSyncEntities(ctx context.Context, entityTypes ...models.EntityType) (models.Response, error) {
	return c.internal.ForceSyncEntitiesSince(ctx, time.Now(), entityTypes...)
}

// ForceSyncEntitiesSince requests complete snapshots of entityTypes together
// with regular changes since lastSync. Pass the server timestamp from the
// previous response as lastSync. An empty entityTypes list behaves like an
// incremental synchronization.
func (c *Client) ForceSyncEntitiesSince(ctx context.Context, lastSync time.Time, entityTypes ...models.EntityType) (models.Response, error) {
	return c.internal.ForceSyncEntitiesSince(ctx, lastSync, entityTypes...)
}

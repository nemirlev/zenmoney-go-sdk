package api

import (
	"context"

	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
)

// Suggest requests merchant, payee, and tag suggestions for transaction. The
// transaction may be partially populated; for example, it may contain only a
// payee. The returned transaction contains the values suggested by ZenMoney.
func (c *Client) Suggest(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return c.internal.Suggest(ctx, transaction)
}

// SuggestBatch requests suggestions for multiple transactions in one API call.
// Each transaction may be partially populated. Results are returned in the same
// order as the input, making this more efficient than repeated Suggest calls.
func (c *Client) SuggestBatch(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, error) {
	return c.internal.SuggestBatch(ctx, transactions)
}

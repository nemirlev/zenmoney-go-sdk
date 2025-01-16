package api

import (
	"context"
	"github.com/nemirlev/zenmoney-go-sdk/models"
)

// Suggest requests suggestions for a single transaction. It can predict merchant, payee,
// and tags based on the provided transaction data. The transaction parameter can be partially
// filled - for example, you can provide only the payee field. The response will contain
// suggested values for the submitted fields.
//
// Example:
//
//	tx := models.Transaction{Payee: "McDonalds"}
//	suggestion, err := client.Suggest(ctx, tx)
//	if err != nil {
//	    // Handle error
//	}
//	// suggestion may contain predicted merchant and tags
func (c *Client) Suggest(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return c.internal.Suggest(ctx, transaction)
}

// SuggestBatch requests suggestions for multiple transactions at once. This is more efficient
// than making multiple single Suggest calls when you need predictions for several transactions.
// Like the Suggest method, transactions can be partially filled, and the response will contain
// predictions only for the provided fields.
//
// Example:
//
//	txs := []models.Transaction{
//	    {Payee: "McDonalds"},
//	    {Payee: "Starbucks"},
//	}
//	suggestions, err := client.SuggestBatch(ctx, txs)
//	if err != nil {
//	    // Handle error
//	}
//	// suggestions will contain predictions for both transactions
//
// SuggestBatch sends a batch suggestion request to the ZenMoney API for multiple transactions.
// It's more efficient than making multiple individual requests when you need suggestions
// for multiple transactions at once.
//
// Parameters:
//   - ctx: Context for the request
//   - transactions: Slice of Transaction objects, each can be partially filled
//
// Returns:
//   - []Transaction: Slice of Transaction objects with suggested values
//   - error: Any error that occurred during the request
//
// The order of suggestions in the response matches the order of input transactions.
func (c *Client) SuggestBatch(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, error) {
	return c.internal.SuggestBatch(ctx, transactions)
}

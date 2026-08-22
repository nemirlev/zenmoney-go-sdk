package models_test

import (
	"encoding/json"
	"testing"

	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
	"github.com/stretchr/testify/require"
)

func TestResponseMatchesCurrentDiffShape(t *testing.T) {
	payload := []byte(`{
		"serverTimestamp": 1787349360,
		"company": [{"fullTitle": null}],
		"user": [{"paidTill": 4797848907, "subscriptionRenewalDate": null}],
		"tag": [{"archive": true}],
		"merchant": [{"mcc": null}],
		"reminder": [{"comment": null}],
		"transaction": [
			{"tag": null},
			{"tag": ["tag-id"]}
		]
	}`)

	var response models.Response
	require.NoError(t, json.Unmarshal(payload, &response))

	require.Equal(t, int64(1787349360), response.ServerTimestamp)
	require.Nil(t, response.Company[0].FullTitle)
	require.Equal(t, int64(4797848907), response.User[0].PaidTill)
	require.Nil(t, response.User[0].SubscriptionRenewalDate)
	require.True(t, response.Tag[0].Archive)
	require.Nil(t, response.Merchant[0].MCC)
	require.Nil(t, response.Reminder[0].Comment)
	require.Nil(t, response.Transaction[0].Tag)
	require.Equal(t, []string{"tag-id"}, response.Transaction[1].Tag)
}

func TestNullableFieldsAcceptValues(t *testing.T) {
	payload := []byte(`{
		"company": [{"fullTitle": "Example Bank"}],
		"merchant": [{"mcc": 5411}],
		"reminder": [{"comment": "Example comment"}],
		"user": [{"subscriptionRenewalDate": 4797848907}]
	}`)

	var response models.Response
	require.NoError(t, json.Unmarshal(payload, &response))

	require.Equal(t, "Example Bank", *response.Company[0].FullTitle)
	require.Equal(t, 5411, *response.Merchant[0].MCC)
	require.Equal(t, "Example comment", *response.Reminder[0].Comment)
	require.Equal(t, int64(4797848907), *response.User[0].SubscriptionRenewalDate)
}

func TestEntityTypeCountry(t *testing.T) {
	require.Equal(t, models.EntityType("country"), models.EntityTypeCountry)
}

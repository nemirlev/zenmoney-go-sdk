package models_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/nemirlev/zenmoney-go-sdk/v2/models"
	"github.com/stretchr/testify/require"
)

func TestResponseMatchesTrackedDiffFixture(t *testing.T) {
	payload, err := os.ReadFile("../example.json")
	require.NoError(t, err)

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()

	var response models.Response
	require.NoError(t, decoder.Decode(&response))
	require.ErrorIs(t, decoder.Decode(&struct{}{}), io.EOF)

	entityTypes := map[string]reflect.Type{
		"instrument":     reflect.TypeOf(models.Instrument{}),
		"country":        reflect.TypeOf(models.Country{}),
		"company":        reflect.TypeOf(models.Company{}),
		"user":           reflect.TypeOf(models.User{}),
		"account":        reflect.TypeOf(models.Account{}),
		"tag":            reflect.TypeOf(models.Tag{}),
		"merchant":       reflect.TypeOf(models.Merchant{}),
		"budget":         reflect.TypeOf(models.Budget{}),
		"reminder":       reflect.TypeOf(models.Reminder{}),
		"reminderMarker": reflect.TypeOf(models.ReminderMarker{}),
		"transaction":    reflect.TypeOf(models.Transaction{}),
		"deletion":       reflect.TypeOf(models.Deletion{}),
	}

	var rawResponse map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(payload, &rawResponse))
	for entityName, rawEntities := range rawResponse {
		if entityName == "serverTimestamp" {
			continue
		}

		modelType, ok := entityTypes[entityName]
		require.True(t, ok, "fixture entity %q has no model type", entityName)

		var entities []map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(rawEntities, &entities), entityName)
		for _, entity := range entities {
			for fieldName, value := range entity {
				if !bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
					continue
				}

				field, found := fieldByJSONName(modelType, fieldName)
				require.True(t, found, "%s.%s has no model field", entityName, fieldName)
				require.Truef(
					t,
					supportsJSONNull(field.Type),
					"%s.%s uses non-nullable Go type %s",
					entityName,
					fieldName,
					field.Type,
				)
			}
		}
	}
}

func TestUnixTimestampFieldsUseInt64(t *testing.T) {
	timestampFields := []struct {
		model reflect.Type
		field string
		want  reflect.Type
	}{
		{reflect.TypeOf(models.Response{}), "ServerTimestamp", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Request{}), "CurrentClientTimestamp", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Request{}), "ServerTimestamp", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Deletion{}), "Stamp", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Instrument{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Company{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.User{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.User{}), "PaidTill", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.User{}), "SubscriptionRenewalDate", reflect.TypeOf((*int64)(nil))},
		{reflect.TypeOf(models.Account{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Tag{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Budget{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Merchant{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Reminder{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.ReminderMarker{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Transaction{}), "Changed", reflect.TypeOf(int64(0))},
		{reflect.TypeOf(models.Transaction{}), "Created", reflect.TypeOf(int64(0))},
	}

	for _, tt := range timestampFields {
		field, ok := tt.model.FieldByName(tt.field)
		require.True(t, ok, "%s.%s is missing", tt.model.Name(), tt.field)
		require.Equal(t, tt.want, field.Type, "%s.%s", tt.model.Name(), tt.field)
	}
}

func fieldByJSONName(modelType reflect.Type, name string) (reflect.StructField, bool) {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if strings.Split(field.Tag.Get("json"), ",")[0] == name {
			return field, true
		}
	}

	return reflect.StructField{}, false
}

func supportsJSONNull(fieldType reflect.Type) bool {
	switch fieldType.Kind() {
	case reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}

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

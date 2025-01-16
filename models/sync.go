package models

// Deletion - information about a deleted object. Some objects, such as Transaction, ReminderMarker, Budget, can be marked as deleted within themselves, but all user objects that have an id can be permanently deleted through deletion. When receiving a deletion object, the receiving party must delete this object on their side.
// Example:
// Suppose the client completely deleted the operation with id '7DE41EB0-3C61-4DB2-BAE8-BDB2A6A46604'. Then in Diff, the following deletion object is transmitted:
//
//	{
//	    //...
//	    deletion: [
//	        {
//	            id: '7DE41EB0-3C61-4DB2-BAE8-BDB2A6A46604',
//	            object: 'transaction',
//	            user: 123456,
//	            stamp: 1490008039
//	        }
//	    ]
//	    //...
//	}
type Deletion struct {
	ID     string `json:"id"`     // Object.id
	Object string `json:"object"` // Object.class
	Stamp  int    `json:"stamp"`
	User   int    `json:"user"`
}

// EntityType represents the type of entity for forceFetch
type EntityType string

const (
	EntityTypeInstrument     EntityType = "instrument"
	EntityTypeCompany        EntityType = "company"
	EntityTypeUser           EntityType = "user"
	EntityTypeAccount        EntityType = "account"
	EntityTypeTag            EntityType = "tag"
	EntityTypeMerchant       EntityType = "merchant"
	EntityTypeBudget         EntityType = "budget"
	EntityTypeReminder       EntityType = "reminder"
	EntityTypeReminderMarker EntityType = "reminderMarker"
	EntityTypeTransaction    EntityType = "transaction"
)

type Response struct {
	ServerTimestamp int              `json:"serverTimestamp"` // Unix timestamp
	Instrument      []Instrument     `json:"instrument,omitempty"`
	Country         []Country        `json:"country,omitempty"`
	Company         []Company        `json:"company,omitempty"`
	User            []User           `json:"user,omitempty"`
	Account         []Account        `json:"account,omitempty"`
	Tag             []Tag            `json:"tag,omitempty"`
	Merchant        []Merchant       `json:"merchant,omitempty"`
	Budget          []Budget         `json:"budget,omitempty"`
	Reminder        []Reminder       `json:"reminder,omitempty"`
	ReminderMarker  []ReminderMarker `json:"reminderMarker,omitempty"`
	Transaction     []Transaction    `json:"transaction,omitempty"`
	Deletion        []Deletion       `json:"deletion,omitempty"`
}

type Request struct {
	CurrentClientTimestamp int              `json:"currentClientTimestamp"` // Unix timestamp
	ServerTimestamp        int              `json:"serverTimestamp"`        // Unix timestamp
	ForceFetch             []EntityType     `json:"forceFetch,omitempty"`
	Instrument             []Instrument     `json:"instrument,omitempty"`
	Country                []Country        `json:"country,omitempty"`
	Company                []Company        `json:"company,omitempty"`
	User                   []User           `json:"user,omitempty"`
	Account                []Account        `json:"account,omitempty"`
	Tag                    []Tag            `json:"tag,omitempty"`
	Merchant               []Merchant       `json:"merchant,omitempty"`
	Budget                 []Budget         `json:"budget,omitempty"`
	Reminder               []Reminder       `json:"reminder,omitempty"`
	ReminderMarker         []ReminderMarker `json:"reminderMarker,omitempty"`
	Transaction            []Transaction    `json:"transaction,omitempty"`
	Deletion               []Deletion       `json:"deletion,omitempty"`
}

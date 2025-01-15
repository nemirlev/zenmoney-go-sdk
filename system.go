package zenmoney

//
// System entities available for read-only access
//

// Instrument - a currency or other financial instrument used by the user to manage their finances.
type Instrument struct {
	// Primary ID of the instrument.
	ID int `json:"id"`

	// Full title of the instrument. For example, United States Dollar, Russian Ruble, Euro.
	Title string `json:"title"`

	// Three-letter code of the instrument. For example, USD, RUB, EUR.
	ShortTitle string `json:"shortTitle"`

	// Symbol of the instrument. For example, $, ₽, €.
	Symbol string `json:"symbol"`

	// Exchange rate of the instrument in rubles.
	Rate float64 `json:"rate"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"`
}

// Company - a bank or other financial organization where accounts can exist.
type Company struct {
	// Primary ID of the company.
	ID int `json:"id"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"`

	// Title of the company.
	Title string `json:"title"`

	// Full title of the company.
	FullTitle string `json:"fullTitle"`

	// Website of the company.
	Www string `json:"www"`

	// Country ID.
	Country int `json:"country"`

	// Indicates if the company is deleted.
	Deleted bool `json:"deleted"`

	// Country code.
	CountryCode string `json:"countryCode"`
}

// User - a ZenMoney user.
type User struct {
	// Primary ID of the user.
	ID int `json:"id"`

	// Country ID of the user.
	Country int `json:"country"`

	// Login name of the user.
	Login string `json:"login"`

	// Parent user ID for family accounting. For the parent user, parent == null.
	Parent *int32 `json:"parent"`

	// Country code of the user. For example, RU, US.
	CountryCode string `json:"countryCode"`

	// Email address of the user.
	Email string `json:"email"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"`

	// Instrument ID. The primary currency used by the user. ID of the Instrument entity.
	Currency int `json:"currency"`

	// Unix timestamp until the user has paid.
	PaidTill int `json:"paidTill"`

	// The day of the month when the user's financial month starts.
	MonthStartDay int `json:"monthStartDay"`

	// Indicates if the forecast feature is enabled for the user.
	IsForecastEnabled bool `json:"isForecastEnabled"`

	// Mode of the plan balance. For example, balance, cashflow.
	PlanBalanceMode string `json:"planBalanceMode"`

	// Settings for the user's plan in JSON array.
	PlanSettings string `json:"planSettings"`

	// Type of subscription the user has. For example, 10YearsSubscription.
	Subscription string `json:"subscription"`

	// Unix timestamp of the subscription renewal date. Null if not applicable.
	SubscriptionRenewalDate *int `json:"subscriptionRenewalDate"`
}

// Country - a country entity with associated information.
type Country struct {
	// Primary ID of the country.
	ID int `json:"id"`

	// Full title of the country. For example, Russia, United States, Germany.
	Title string `json:"title"`

	// Instrument.id. The primary currency used in the country.
	Currency int `json:"currency"`

	// Domain code of the country. For example, ru, us, de.
	Domain string `json:"domain"`
}

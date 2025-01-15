package zenmoney

//
// User entities available for reading, creating, updating, and deleting
//

// Account - user account
type Account struct {
	// Primary ID of the account. UUID.
	ID string `json:"id"`

	// User ID associated with the account. See User.ID
	User int `json:"user"`

	// Instrument ID associated with the account. See Instrument.ID
	Instrument *int32 `json:"instrument"`

	// Type of the account. For example, ccard - credit card, checking - bank account, loan, deposit, cash, debt.
	Type string `json:"type"`

	// Role ID associated with the account. See User.ID
	Role *int32 `json:"role"`

	// Indicates if the account is private.
	Private bool `json:"private"`

	// Indicates if the account is a savings account.
	Savings *bool `json:"savings"`

	// Title of the account. For example, 'My main Sberbank', 'Alfa-Bank Premium', 'Cash'.
	Title string `json:"title"`

	// Indicates if the account is included in the balance. If true, its balance is included in the overall balance and expenses and income are accounted for in reports.
	InBalance bool `json:"inBalance"`

	// Credit limit if the account type is a credit card or bank account.
	CreditLimit *float64 `json:"creditLimit"`

	// Balance of the account at the time of opening. If the account type is a loan, the principal amount is specified.
	StartBalance *float64 `json:"startBalance"`

	// Current balance of the account.
	Balance *float64 `json:"balance"`

	// Company ID associated with the account. Company.ID
	Company *int32 `json:"company"`

	// Indicates if the account is archived.
	Archive bool `json:"archive"`

	// If true, the ZenMoney app will adjust the account balance to its value in the SMS when recognizing SMS.
	EnableCorrection bool `json:"enableCorrection"`

	// Type of balance correction. For example, 'request'.
	BalanceCorrectionType string `json:"balanceCorrectionType"`

	// Date of opening the deposit/loan.
	StartDate *string `json:"startDate"`

	// For deposits - whether interest capitalization is enabled. For loans - whether the loan is annuity.
	Capitalization *bool `json:"capitalization"`

	// Interest rate on the account (in percent).
	Percent *float64 `json:"percent"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"`

	// Array of bank account numbers. Usually, the last 4 digits of the account number and the last 4 digits of the bank cards linked to the account.
	SyncID []string `json:"syncID"`

	// Indicates if SMS recognition is enabled for the account.
	EnableSMS bool `json:"enableSMS"`

	// Term of the loan/deposit in intervals of endDateOffsetInterval starting from the opening date.
	EndDateOffset *int32 `json:"endDateOffset"`

	// Interval between payments. For example, 'day', 'week', 'month', 'year'.
	EndDateOffsetInterval *string `json:"endDateOffsetInterval"`

	// Number of payoff intervals after which payments are made, starting from the next date startDate + payoffInterval.
	PayoffStep *int32 `json:"payoffStep"`

	// Interval for payments. For example, 'month', 'year'.
	PayoffInterval *string `json:"payoffInterval"`
}

// Tag - operation category
type Tag struct {
	// Primary ID of the tag.
	ID string `json:"id"` // UUID

	// User ID associated with the tag.
	User int `json:"user"` // User.id

	// Unix timestamp of the last change.
	Changed int `json:"changed"` // Unix timestamp

	// ID of the tag icon.
	Icon *string `json:"icon"` // Id иконки категории

	// Indicates if the tag is included in the income budget calculation.
	BudgetIncome bool `json:"budgetIncome"` // Включена ли категория в расчёт дохода в бюджете

	// Indicates if the tag is included in the outcome budget calculation.
	BudgetOutcome bool `json:"budgetOutcome"` // Включена ли категория в расчёт расхода в бюджете

	// Indicates if the expenses for this tag are mandatory. If null, they are also considered mandatory.
	Required *bool `json:"required"` // Являются ли расходы по данной категории обязательными

	// Color of the tag icon as a number. Calculated by alpha, red, green, blue (0 <= 255).
	Color *int64 `json:"color"` // Цвет иконки категории

	// URL of the picture for the tag.
	Picture *string `json:"picture"` // Ссылка на картинку для данной категории

	// Title of the tag.
	Title string `json:"title"`

	// Indicates if the tag is for income.
	ShowIncome bool `json:"showIncome"` // Является ли категория доходной

	// Indicates if the tag is for outcome.
	ShowOutcome bool `json:"showOutcome"` // Является ли категория расходной

	// ID of the parent tag. The parent category. Nested level should not exceed 1.
	Parent *string `json:"parent"` // Tag.id

	// Static ID of the tag.
	StaticID string `json:"staticId"` // Статический ID категории
}

// Budget - user budget
// You can delete the budget for the current month by removing the lock and setting the corresponding amount to 0.
type Budget struct {
	// User ID associated with the budget.
	User int `json:"user"` // User.id

	// Unix timestamp of the last change.
	Changed int `json:"changed"` // Unix timestamp

	// Start date of the budget month.
	Date string `json:"date"` // 'yyyy-MM-dd'

	// Tag ID associated with the budget.
	Tag *string `json:"tag"` // Tag.id

	// Income budget amount.
	Income float64 `json:"income"` // Доходный бюджет

	// Outcome budget amount.
	Outcome float64 `json:"outcome"` // Расходный бюджет

	// Indicates if the income budget is locked.
	IncomeLock bool `json:"incomeLock"` // Блокировка доходного бюджета

	// Indicates if the outcome budget is locked.
	OutcomeLock bool `json:"outcomeLock"` // Блокировка расходного бюджета

	// Indicates if the income forecast is enabled.
	IsIncomeForecast bool `json:"isIncomeForecast"` // Включен ли прогноз доходов

	// Indicates if the outcome forecast is enabled.
	IsOutcomeForecast bool `json:"isOutcomeForecast"` // Включен ли прогноз расходов
}

// Merchant - operation counterparty (seller, supplier, payer)
type Merchant struct {
	// Primary ID of the merchant.
	ID string `json:"id"` // UUID

	// User ID associated with the merchant.
	User int `json:"user"` // User.id

	// Title of the merchant.
	Title string `json:"title"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"` // Unix timestamp
}

// Reminder - a reminder for a financial operation
// Example:
// A Reminder with parameters interval: 'day', step: 7, points: [0, 2, 4], startDate: '2017-03-08', endDate: null
// means that operations need to be repeated every week starting from 2017-03-08 on Wednesdays, Fridays, and Sundays.
// Because 2017-03-08 is a Wednesday, point 0 is Wednesday, point 2 is Friday, and point 4 is Sunday.
// Every week - because the step is 7 days.
type Reminder struct {
	// Primary ID of the reminder. UUID.
	ID string `json:"id"`

	// User.ID associated with the reminder.
	User int `json:"user"`

	// Income amount for the reminder.
	Income float64 `json:"income"`

	// Outcome amount for the reminder.
	Outcome float64 `json:"outcome"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"`

	// Instrument.ID for the income.
	IncomeInstrument int `json:"incomeInstrument"`

	// Instrument.ID for the outcome.
	OutcomeInstrument int `json:"outcomeInstrument"`

	// Step interval for the reminder.
	Step int `json:"step"`

	// Points in time for the reminder.
	Points []int `json:"points"`

	// Tags associated with the reminder. Array of Tag.ID.
	Tag []string `json:"tag"`

	// Start date of the reminder. Format 'yyyy-MM-dd'.
	StartDate string `json:"startDate"`

	// End date of the reminder. Format 'yyyy-MM-dd'.
	EndDate *string `json:"endDate"`

	// Indicates if notifications are enabled for the reminder.
	Notify bool `json:"notify"`

	// Interval for the reminder.
	Interval *string `json:"interval"`

	// Account.ID for the income.
	IncomeAccount string `json:"incomeAccount"`

	// Account.ID for the outcome.
	OutcomeAccount string `json:"outcomeAccount"`

	// Comment for the reminder.
	Comment string `json:"comment"`

	// Payee.ID associated with the reminder.
	Payee *string `json:"payee"`

	// Merchant.ID associated with the reminder.
	Merchant *string `json:"merchant"`
}

// ReminderMarker - a marker for a reminder
type ReminderMarker struct {
	// Primary ID of the reminder marker. UUID.
	ID string `json:"id"`

	// User ID associated with the reminder marker. See User.ID
	User int `json:"user"`

	// Date of the reminder marker. In format 'yyyy-MM-dd'.
	Date string `json:"date"`

	// Income amount for the reminder marker.
	Income float64 `json:"income"`

	// Outcome amount for the reminder marker.
	Outcome float64 `json:"outcome"`

	// Unix timestamp of the last change. In Unix timestamp format.
	Changed int `json:"changed"`

	// Instrument.ID for the income.
	IncomeInstrument int `json:"incomeInstrument"`

	// Instrument.ID for the outcome.
	OutcomeInstrument int `json:"outcomeInstrument"`

	// State of the reminder marker.
	State string `json:"state"`

	// Indicates if the reminder marker is a forecast.
	IsForecast bool `json:"isForecast"`

	// Reminder.ID associated with the reminder marker.
	Reminder string `json:"reminder"`

	// Account.ID for the income.
	IncomeAccount string `json:"incomeAccount"`

	// Account.ID for the outcome.
	OutcomeAccount string `json:"outcomeAccount"`

	// Comment for the reminder marker.
	Comment string `json:"comment"`

	// Payee associated with the reminder marker.
	Payee *string `json:"payee"`

	// Merchant associated with the reminder marker.
	Merchant *string `json:"merchant"`

	// Indicates if notifications are enabled for the reminder marker.
	Notify bool `json:"notify"`

	// Tag.ID associated with the reminder marker.
	Tag []string `json:"tag"`
}

// Transaction - a financial transaction
//
//	income: incomeAccount=outcomeAccount && income > 0
//	outcome: incomeAccount=outcomeAccount && outcome > 0
//	debt income: type(incomeAccount) == "debt"
//	debt outcome: type(outcomeAccount) == "debt"
//	transfer: other
//
// Examples: https://github.com/zenmoney/ZenPlugins/wiki/ZenMoney-API#transaction
type Transaction struct {
	// Primary ID of the transaction. UUID.
	ID string `json:"id"`

	// User.ID associated with the transaction.
	User int `json:"user"`

	// Date of the transaction. Format 'yyyy-MM-dd'.
	Date string `json:"date"`

	// Income amount for the transaction.
	Income float64 `json:"income"`

	// Outcome amount for the transaction.
	Outcome float64 `json:"outcome"`

	// Unix timestamp of the last change.
	Changed int `json:"changed"`

	// Instrument.ID for the income.
	IncomeInstrument int `json:"incomeInstrument"`

	// Instrument.ID for the outcome.
	OutcomeInstrument int `json:"outcomeInstrument"`

	// Unix timestamp of the creation of the transaction.
	Created int `json:"created"` // Unix timestamp

	// Original payee of the transaction.
	OriginalPayee string `json:"originalPayee"`

	// Indicates if the transaction is deleted.
	Deleted bool `json:"deleted"`

	// Indicates if the transaction is viewed.
	Viewed bool `json:"viewed"`

	// Indicates if the transaction is on hold.
	Hold bool `json:"hold"`

	// QR code associated with the transaction.
	QRCode *string `json:"qrCode"`

	// Source of the transaction.
	Source string `json:"source"`

	// Account.ID for the income.
	IncomeAccount string `json:"incomeAccount"`

	// Account.ID for the outcome.
	OutcomeAccount *string `json:"outcomeAccount"`

	// Tags associated with the transaction. Array of Tag.ID.
	Tag []string `json:"tag"`

	// Comment for the transaction.
	Comment *string `json:"comment"`

	// Payee of the transaction.
	Payee string `json:"payee"`

	// Operational income amount.
	OpIncome float64 `json:"opIncome"`

	// Operational outcome amount.
	OpOutcome float64 `json:"opOutcome"`

	// Instrument.ID for the operational income.
	OpIncomeInstrument *int `json:"opIncomeInstrument"`

	// Instrument.ID for the operational outcome.
	OpOutcomeInstrument *int `json:"opOutcomeInstrument"`

	// Latitude of the transaction location.
	Latitude *float64 `json:"latitude"` // Широта

	// Longitude of the transaction location.
	Longitude *float64 `json:"longitude"` // Долгота

	// Merchant.ID associated with the transaction.
	Merchant *string `json:"merchant"`

	// Bank ID for the income.
	IncomeBankID *string `json:"incomeBankID"`

	// Bank ID for the outcome.
	OutcomeBankID *string `json:"outcomeBankID"`

	// Reminder.ID associated with the transaction.
	ReminderMarker *string `json:"reminderMarker"`
}

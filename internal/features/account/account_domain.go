package account

import "time"

type AccountCurrency string

const (
	CurrencyTHB AccountCurrency = "THB"
	CurrencyUSD AccountCurrency = "USD"
)

type Account struct {
	ID        int64           `json:"id"`
	OwnerID   int64           `json:"owner_id"`
	Balance   int64           `json:"balance"`
	Currency  AccountCurrency `json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

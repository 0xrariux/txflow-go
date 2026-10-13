package types

import (
	"github.com/0xrariux/txflow-go/decimal"
	"time"
)

type Balance struct {
	Asset     string          `json:"asset"`
	Total     decimal.Decimal `json:"total"`
	Available decimal.Decimal `json:"available"`
}
type Position struct {
	Symbol        string          `json:"symbol"`
	Quantity      decimal.Decimal `json:"quantity"`
	EntryPrice    decimal.Decimal `json:"entry_price"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	MarginMode    string          `json:"margin_mode"`
}
type AccountState struct {
	Address    string           `json:"address"`
	Balances   []Balance        `json:"balances"`
	Positions  []Position       `json:"positions"`
	Equity     *decimal.Decimal `json:"equity,omitempty"`
	MarginUsed *decimal.Decimal `json:"margin_used,omitempty"`
}
type AccountOrder struct {
	ID            string          `json:"id"`
	ClientOrderID string          `json:"client_order_id,omitempty"`
	Symbol        string          `json:"symbol"`
	Side          Side            `json:"side"`
	Price         decimal.Decimal `json:"price"`
	Quantity      decimal.Decimal `json:"quantity"`
	Status        string          `json:"status"`
}
type Fill struct {
	Trade   Trade           `json:"trade"`
	OrderID string          `json:"order_id"`
	Fee     decimal.Decimal `json:"fee"`
}
type FundingPayment struct {
	Symbol    string          `json:"symbol"`
	Amount    decimal.Decimal `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

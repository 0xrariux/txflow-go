// Package types defines proposed SDK domain models. Their JSON representation
// is the local demo schema, NOT a verified TxFlow API wire format.
package types

import (
	"github.com/0xrariux/txflow-go/decimal"
	"time"
)

type MarketKind string

const (
	Spot       MarketKind = "spot"
	Perpetual  MarketKind = "perpetual"
	Prediction MarketKind = "prediction"
)

type Market struct {
	Symbol            string          `json:"symbol"`
	Kind              MarketKind      `json:"kind"`
	BaseAsset         string          `json:"base_asset"`
	QuoteAsset        string          `json:"quote_asset"`
	PriceIncrement    decimal.Decimal `json:"price_increment"`
	QuantityIncrement decimal.Decimal `json:"quantity_increment"`
}
type Asset struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}
type Level struct {
	Price decimal.Decimal `json:"price"`
	Size  decimal.Decimal `json:"size"`
}
type OrderBookRequest struct {
	Symbol string
	Depth  int
}
type OrderBook struct {
	Symbol    string    `json:"symbol"`
	Bids      []Level   `json:"bids"`
	Asks      []Level   `json:"asks"`
	Timestamp time.Time `json:"timestamp"`
	// Sequence is populated only by an adapter with a known sequence contract.
	Sequence *uint64 `json:"sequence,omitempty"`
}
type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

type Trade struct {
	ID        string          `json:"id"`
	Symbol    string          `json:"symbol"`
	Side      Side            `json:"side"`
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	Timestamp time.Time       `json:"timestamp"`
}
type Ticker struct {
	Symbol string          `json:"symbol"`
	Last   decimal.Decimal `json:"last"`
	// Optional fields distinguish unavailable values from a genuine zero.
	Mark         *decimal.Decimal `json:"mark,omitempty"`
	Oracle       *decimal.Decimal `json:"oracle,omitempty"`
	FundingRate  *decimal.Decimal `json:"funding_rate,omitempty"`
	OpenInterest *decimal.Decimal `json:"open_interest,omitempty"`
	Timestamp    time.Time        `json:"timestamp"`
}

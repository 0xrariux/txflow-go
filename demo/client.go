// Package demo explicitly provides authored synthetic data, never live TxFlow
// data. It performs no network I/O and defines no TxFlow wire schema.
package demo

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	txflow "github.com/0xrariux/txflow-go"
	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

//go:embed testdata/markets.json
var marketData []byte

type dataset struct {
	Markets []types.Market    `json:"markets"`
	Assets  []types.Asset     `json:"assets"`
	Books   []types.OrderBook `json:"books"`
	Trades  []types.Trade     `json:"trades"`
	Tickers []types.Ticker    `json:"tickers"`
}
type reader struct{ data dataset }

func NewClient() (*txflow.Client, error) {
	var data dataset
	if err := json.Unmarshal(marketData, &data); err != nil {
		return nil, fmt.Errorf("demo: decode synthetic fixture: %w", err)
	}
	var accounts accountDataset
	if err := json.Unmarshal(accountData, &accounts); err != nil {
		return nil, fmt.Errorf("demo: account fixture: %w", err)
	}
	connector, err := newConnector()
	if err != nil {
		return nil, err
	}
	return txflow.NewClient(txflow.WithMarketReader(&reader{data: data}), txflow.WithAccountReader(&accountReader{data: accounts}), txflow.WithStreamConnector(connector))
}
func (r *reader) List(ctx context.Context) ([]types.Market, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return append([]types.Market(nil), r.data.Markets...), nil
}
func (r *reader) Assets(ctx context.Context) ([]types.Asset, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return append([]types.Asset(nil), r.data.Assets...), nil
}
func (r *reader) OrderBook(ctx context.Context, req types.OrderBookRequest) (types.OrderBook, error) {
	if err := ctx.Err(); err != nil {
		return types.OrderBook{}, err
	}
	for _, book := range r.data.Books {
		if book.Symbol == req.Symbol {
			n := req.Depth
			if n > len(book.Bids) {
				n = len(book.Bids)
			}
			book.Bids = append([]types.Level(nil), book.Bids[:n]...)
			n = req.Depth
			if n > len(book.Asks) {
				n = len(book.Asks)
			}
			book.Asks = append([]types.Level(nil), book.Asks[:n]...)
			if book.Sequence != nil {
				s := *book.Sequence
				book.Sequence = &s
			}
			return book, nil
		}
	}
	return types.OrderBook{}, sdkerr.ErrNotFound
}
func (r *reader) Trades(ctx context.Context, symbol string) ([]types.Trade, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var out []types.Trade
	for _, t := range r.data.Trades {
		if t.Symbol == symbol {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return nil, sdkerr.ErrNotFound
	}
	return out, nil
}
func cloneDecimal(p *types.Ticker) types.Ticker {
	out := *p
	if p.Mark != nil {
		v := *p.Mark
		out.Mark = &v
	}
	if p.Oracle != nil {
		v := *p.Oracle
		out.Oracle = &v
	}
	if p.FundingRate != nil {
		v := *p.FundingRate
		out.FundingRate = &v
	}
	if p.OpenInterest != nil {
		v := *p.OpenInterest
		out.OpenInterest = &v
	}
	return out
}
func (r *reader) Ticker(ctx context.Context, symbol string) (types.Ticker, error) {
	if err := ctx.Err(); err != nil {
		return types.Ticker{}, err
	}
	for _, ticker := range r.data.Tickers {
		if ticker.Symbol == symbol {
			return cloneDecimal(&ticker), nil
		}
	}
	return types.Ticker{}, sdkerr.ErrNotFound
}

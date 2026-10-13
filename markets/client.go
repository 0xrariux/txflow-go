// Package markets supplies proposed typed read operations via an explicit
// backend. There is no built-in TxFlow endpoint implementation.
package markets

import (
	"context"
	"fmt"
	"strings"

	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

// Reader must be concurrency-safe; returned values are caller-owned.
// A future verified adapter maps wire models into these domain models.
type Reader interface {
	List(context.Context) ([]types.Market, error)
	Assets(context.Context) ([]types.Asset, error)
	OrderBook(context.Context, types.OrderBookRequest) (types.OrderBook, error)
	Trades(context.Context, string) ([]types.Trade, error)
	Ticker(context.Context, string) (types.Ticker, error)
}
type Client struct{ reader Reader }

func New(r Reader) *Client { return &Client{reader: r} }
func (c *Client) ready(ctx context.Context) error {
	if ctx == nil {
		return sdkerr.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if c.reader == nil {
		return sdkerr.ErrPendingConfirmation
	}
	return nil
}
func ValidSymbol(s string) bool {
	return len(s) > 0 && len(s) <= 128 && strings.TrimSpace(s) == s && !strings.ContainsAny(s, "\r\n\t")
}
func (c *Client) List(ctx context.Context) ([]types.Market, error) {
	if err := c.ready(ctx); err != nil {
		return nil, err
	}
	out, err := c.reader.List(ctx)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, m := range out {
		if !ValidSymbol(m.Symbol) || seen[m.Symbol] || m.BaseAsset == "" || m.QuoteAsset == "" || m.PriceIncrement.Sign() <= 0 || m.QuantityIncrement.Sign() <= 0 {
			return nil, fmt.Errorf("%w: invalid market metadata", sdkerr.ErrMalformedResponse)
		}
		if m.Kind != types.Spot && m.Kind != types.Perpetual && m.Kind != types.Prediction {
			return nil, sdkerr.ErrMalformedResponse
		}
		seen[m.Symbol] = true
	}
	return out, nil
}
func (c *Client) Metadata(ctx context.Context, symbol string) (types.Market, error) {
	if !ValidSymbol(symbol) {
		return types.Market{}, sdkerr.ErrInvalidRequest
	}
	out, err := c.List(ctx)
	if err != nil {
		return types.Market{}, err
	}
	for _, m := range out {
		if m.Symbol == symbol {
			return m, nil
		}
	}
	return types.Market{}, sdkerr.ErrNotFound
}
func (c *Client) Assets(ctx context.Context) ([]types.Asset, error) {
	if err := c.ready(ctx); err != nil {
		return nil, err
	}
	out, err := c.reader.Assets(ctx)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, a := range out {
		if a.Symbol == "" || a.Name == "" || seen[a.Symbol] {
			return nil, sdkerr.ErrMalformedResponse
		}
		seen[a.Symbol] = true
	}
	return out, nil
}
func (c *Client) OrderBook(ctx context.Context, req types.OrderBookRequest) (types.OrderBook, error) {
	if !ValidSymbol(req.Symbol) || req.Depth < 0 || req.Depth > 1000 {
		return types.OrderBook{}, sdkerr.ErrInvalidRequest
	}
	if err := c.ready(ctx); err != nil {
		return types.OrderBook{}, err
	}
	if req.Depth == 0 {
		req.Depth = 20
	}
	out, err := c.reader.OrderBook(ctx, req)
	if err != nil {
		return types.OrderBook{}, err
	}
	if out.Symbol != req.Symbol || len(out.Bids) > req.Depth || len(out.Asks) > req.Depth {
		return types.OrderBook{}, sdkerr.ErrMalformedResponse
	}
	for _, side := range [][]types.Level{out.Bids, out.Asks} {
		seen := map[string]bool{}
		for _, l := range side {
			if l.Price.Sign() <= 0 || l.Size.Sign() <= 0 || seen[l.Price.String()] {
				return types.OrderBook{}, sdkerr.ErrMalformedResponse
			}
			seen[l.Price.String()] = true
		}
	}

	var bid, ask *types.Level
	for i := range out.Bids {
		if bid == nil || out.Bids[i].Price.Cmp(bid.Price) > 0 {
			bid = &out.Bids[i]
		}
	}
	for i := range out.Asks {
		if ask == nil || out.Asks[i].Price.Cmp(ask.Price) < 0 {
			ask = &out.Asks[i]
		}
	}
	if bid != nil && ask != nil && bid.Price.Cmp(ask.Price) >= 0 {
		return types.OrderBook{}, sdkerr.ErrMalformedResponse
	}
	return out, nil
}
func (c *Client) Trades(ctx context.Context, symbol string) ([]types.Trade, error) {
	if !ValidSymbol(symbol) {
		return nil, sdkerr.ErrInvalidRequest
	}
	if err := c.ready(ctx); err != nil {
		return nil, err
	}
	out, err := c.reader.Trades(ctx, symbol)
	if err != nil {
		return nil, err
	}
	for _, t := range out {
		if t.Symbol != symbol || t.ID == "" || t.Price.Sign() <= 0 || t.Size.Sign() <= 0 || (t.Side != types.Buy && t.Side != types.Sell) {
			return nil, sdkerr.ErrMalformedResponse
		}
	}
	return out, nil
}
func (c *Client) Ticker(ctx context.Context, symbol string) (types.Ticker, error) {
	if !ValidSymbol(symbol) {
		return types.Ticker{}, sdkerr.ErrInvalidRequest
	}
	if err := c.ready(ctx); err != nil {
		return types.Ticker{}, err
	}
	out, err := c.reader.Ticker(ctx, symbol)
	if err != nil {
		return types.Ticker{}, err
	}
	if out.Symbol != symbol || out.Last.Sign() <= 0 || (out.Mark != nil && out.Mark.Sign() <= 0) || (out.Oracle != nil && out.Oracle.Sign() <= 0) || (out.OpenInterest != nil && out.OpenInterest.Sign() < 0) {
		return types.Ticker{}, sdkerr.ErrMalformedResponse
	}
	return out, nil
}

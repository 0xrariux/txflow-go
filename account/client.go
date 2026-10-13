// Package account defines proposed query-only domain operations. No private
// key is needed. There is no built-in TxFlow endpoint implementation.
package account

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

type Reader interface {
	State(context.Context, string) (types.AccountState, error)
	OpenOrders(context.Context, string) ([]types.AccountOrder, error)
	HistoricalOrders(context.Context, string) ([]types.AccountOrder, error)
	Fills(context.Context, string) ([]types.Fill, error)
	FundingHistory(context.Context, string) ([]types.FundingPayment, error)
}
type Client struct{ reader Reader }

func New(r Reader) *Client { return &Client{reader: r} }

// ValidAddress checks the proposed EVM address shape, not TxFlow eligibility.
func ValidAddress(s string) bool {
	if len(s) != 42 || !strings.HasPrefix(s, "0x") {
		return false
	}
	_, err := hex.DecodeString(s[2:])
	return err == nil
}
func (c *Client) ready(ctx context.Context, address string) error {
	if ctx == nil || !ValidAddress(address) {
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
func (c *Client) State(ctx context.Context, address string) (types.AccountState, error) {
	if err := c.ready(ctx, address); err != nil {
		return types.AccountState{}, err
	}
	out, err := c.reader.State(ctx, address)
	if err != nil {
		return types.AccountState{}, err
	}
	if !strings.EqualFold(out.Address, address) {
		return types.AccountState{}, sdkerr.ErrMalformedResponse
	}
	for _, b := range out.Balances {
		if b.Asset == "" {
			return types.AccountState{}, sdkerr.ErrMalformedResponse
		}
	}
	for _, p := range out.Positions {
		if p.Symbol == "" || p.EntryPrice.Sign() < 0 {
			return types.AccountState{}, sdkerr.ErrMalformedResponse
		}
	}
	return out, nil
}
func validateOrders(out []types.AccountOrder) error {
	for _, o := range out {
		if o.ID == "" || o.Symbol == "" || o.Quantity.Sign() <= 0 || o.Price.Sign() < 0 || (o.Side != types.Buy && o.Side != types.Sell) {
			return sdkerr.ErrMalformedResponse
		}
	}
	return nil
}
func (c *Client) OpenOrders(ctx context.Context, address string) ([]types.AccountOrder, error) {
	if err := c.ready(ctx, address); err != nil {
		return nil, err
	}
	out, err := c.reader.OpenOrders(ctx, address)
	if err != nil {
		return nil, err
	}
	if err := validateOrders(out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *Client) HistoricalOrders(ctx context.Context, address string) ([]types.AccountOrder, error) {
	if err := c.ready(ctx, address); err != nil {
		return nil, err
	}
	out, err := c.reader.HistoricalOrders(ctx, address)
	if err != nil {
		return nil, err
	}
	if err := validateOrders(out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *Client) Fills(ctx context.Context, address string) ([]types.Fill, error) {
	if err := c.ready(ctx, address); err != nil {
		return nil, err
	}
	out, err := c.reader.Fills(ctx, address)
	if err != nil {
		return nil, err
	}
	for _, f := range out {
		if f.OrderID == "" || f.Trade.ID == "" || f.Trade.Symbol == "" || f.Trade.Price.Sign() <= 0 || f.Trade.Size.Sign() <= 0 || (f.Trade.Side != types.Buy && f.Trade.Side != types.Sell) {
			return nil, sdkerr.ErrMalformedResponse
		}
	}
	return out, nil
}
func (c *Client) FundingHistory(ctx context.Context, address string) ([]types.FundingPayment, error) {
	if err := c.ready(ctx, address); err != nil {
		return nil, err
	}
	out, err := c.reader.FundingHistory(ctx, address)
	if err != nil {
		return nil, err
	}
	for _, f := range out {
		if f.Symbol == "" || f.Timestamp.IsZero() {
			return nil, sdkerr.ErrMalformedResponse
		}
	}
	return out, nil
}

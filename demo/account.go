package demo

import (
	"context"
	_ "embed"

	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

// Address is the fixture's synthetic public address. It is not a user account.
const Address = "0x0000000000000000000000000000000000000000"

//go:embed testdata/account.json
var accountData []byte

type accountDataset struct {
	State            types.AccountState     `json:"state"`
	OpenOrders       []types.AccountOrder   `json:"open_orders"`
	HistoricalOrders []types.AccountOrder   `json:"historical_orders"`
	Fills            []types.Fill           `json:"fills"`
	Funding          []types.FundingPayment `json:"funding"`
}
type accountReader struct{ data accountDataset }

func (r *accountReader) check(ctx context.Context, address string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if address != Address {
		return sdkerr.ErrNotFound
	}
	return nil
}
func (r *accountReader) State(ctx context.Context, address string) (types.AccountState, error) {
	if err := r.check(ctx, address); err != nil {
		return types.AccountState{}, err
	}
	out := r.data.State
	out.Balances = append([]types.Balance(nil), out.Balances...)
	out.Positions = append([]types.Position(nil), out.Positions...)
	if out.Equity != nil {
		v := *out.Equity
		out.Equity = &v
	}
	if out.MarginUsed != nil {
		v := *out.MarginUsed
		out.MarginUsed = &v
	}
	return out, nil
}
func (r *accountReader) OpenOrders(ctx context.Context, address string) ([]types.AccountOrder, error) {
	if err := r.check(ctx, address); err != nil {
		return nil, err
	}
	return append([]types.AccountOrder(nil), r.data.OpenOrders...), nil
}
func (r *accountReader) HistoricalOrders(ctx context.Context, address string) ([]types.AccountOrder, error) {
	if err := r.check(ctx, address); err != nil {
		return nil, err
	}
	return append([]types.AccountOrder(nil), r.data.HistoricalOrders...), nil
}
func (r *accountReader) Fills(ctx context.Context, address string) ([]types.Fill, error) {
	if err := r.check(ctx, address); err != nil {
		return nil, err
	}
	return append([]types.Fill(nil), r.data.Fills...), nil
}
func (r *accountReader) FundingHistory(ctx context.Context, address string) ([]types.FundingPayment, error) {
	if err := r.check(ctx, address); err != nil {
		return nil, err
	}
	return append([]types.FundingPayment(nil), r.data.Funding...), nil
}

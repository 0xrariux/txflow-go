package account

import (
	"context"
	"errors"
	"testing"

	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

const address = "0x0000000000000000000000000000000000000000"

type corrupt struct{}

func (corrupt) State(context.Context, string) (types.AccountState, error) {
	return types.AccountState{Address: "wrong"}, nil
}
func (corrupt) OpenOrders(context.Context, string) ([]types.AccountOrder, error) {
	return []types.AccountOrder{{}}, nil
}
func (corrupt) HistoricalOrders(context.Context, string) ([]types.AccountOrder, error) {
	return []types.AccountOrder{{}}, nil
}
func (corrupt) Fills(context.Context, string) ([]types.Fill, error) { return []types.Fill{{}}, nil }
func (corrupt) FundingHistory(context.Context, string) ([]types.FundingPayment, error) {
	return []types.FundingPayment{{}}, nil
}
func TestAccountValidation(t *testing.T) {
	c := New(corrupt{})
	ctx := context.Background()
	_, e1 := c.State(ctx, address)
	_, e2 := c.OpenOrders(ctx, address)
	_, e3 := c.HistoricalOrders(ctx, address)
	_, e4 := c.Fills(ctx, address)
	_, e5 := c.FundingHistory(ctx, address)
	for _, e := range []error{e1, e2, e3, e4, e5} {
		if !errors.Is(e, sdkerr.ErrMalformedResponse) {
			t.Fatal(e)
		}
	}
	if _, err := c.State(nil, address); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := c.State(canceled, address); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	for _, s := range []string{"", "0x123", "0xgg00000000000000000000000000000000000000"} {
		if ValidAddress(s) {
			t.Fatal(s)
		}
	}
}

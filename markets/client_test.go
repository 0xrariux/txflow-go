package markets

import (
	"context"
	"errors"
	"testing"

	"github.com/0xrariux/txflow-go/decimal"
	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

type corruptReader struct{}

func (corruptReader) List(context.Context) ([]types.Market, error) {
	return []types.Market{{Symbol: "BTC", Kind: types.Perpetual}}, nil
}
func (corruptReader) Assets(context.Context) ([]types.Asset, error) { return []types.Asset{{}}, nil }
func (corruptReader) OrderBook(context.Context, types.OrderBookRequest) (types.OrderBook, error) {
	return types.OrderBook{Symbol: "BTC", Bids: []types.Level{{Price: decimal.MustParse("1"), Size: decimal.MustParse("-1")}}}, nil
}
func (corruptReader) Trades(context.Context, string) ([]types.Trade, error) {
	return []types.Trade{{}}, nil
}
func (corruptReader) Ticker(context.Context, string) (types.Ticker, error) {
	return types.Ticker{}, nil
}
func TestMalformedAdapterOutput(t *testing.T) {
	c := New(corruptReader{})
	ctx := context.Background()
	_, e1 := c.List(ctx)
	_, e2 := c.Assets(ctx)
	_, e3 := c.OrderBook(ctx, types.OrderBookRequest{Symbol: "BTC"})
	_, e4 := c.Trades(ctx, "BTC")
	_, e5 := c.Ticker(ctx, "BTC")
	for _, err := range []error{e1, e2, e3, e4, e5} {
		if !errors.Is(err, sdkerr.ErrMalformedResponse) {
			t.Fatal(err)
		}
	}
}

type badBookReader struct{ corruptReader }

func (badBookReader) OrderBook(context.Context, types.OrderBookRequest) (types.OrderBook, error) {
	return types.OrderBook{Symbol: "BTC", Bids: []types.Level{{Price: decimal.MustParse("2"), Size: decimal.MustParse("1")}}, Asks: []types.Level{{Price: decimal.MustParse("1"), Size: decimal.MustParse("1")}}}, nil
}
func TestCrossedAdapterBookRejected(t *testing.T) {
	_, err := New(badBookReader{}).OrderBook(context.Background(), types.OrderBookRequest{Symbol: "BTC"})
	if !errors.Is(err, sdkerr.ErrMalformedResponse) {
		t.Fatal(err)
	}
}

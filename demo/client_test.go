package demo_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	txflow "github.com/0xrariux/txflow-go"
	"github.com/0xrariux/txflow-go/decimal"
	"github.com/0xrariux/txflow-go/demo"
	"github.com/0xrariux/txflow-go/types"
)

func TestSyntheticMarketFixtureAndIsolation(t *testing.T) {
	c, err := demo.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	rows, err := c.Markets.List(ctx)
	if err != nil || len(rows) != 2 {
		t.Fatal(rows, err)
	}
	rows[0].Symbol = "changed"
	m, err := c.Markets.Metadata(ctx, "BTC-USDC")
	if err != nil || m.Symbol != "BTC-USDC" {
		t.Fatal(m, err)
	}
	assets, err := c.Markets.Assets(ctx)
	if err != nil || len(assets) != 3 {
		t.Fatal(assets, err)
	}
	book, err := c.Markets.OrderBook(ctx, types.OrderBookRequest{Symbol: "BTC-USDC", Depth: 1})
	if err != nil || len(book.Bids) != 1 {
		t.Fatal(book, err)
	}
	book.Bids[0].Price = decimal.MustParse("1")
	*book.Sequence = 0
	copy, err := c.Markets.OrderBook(ctx, types.OrderBookRequest{Symbol: "BTC-USDC", Depth: 1})
	if err != nil || copy.Bids[0].Price.String() != "69999.9" || *copy.Sequence != 100 {
		t.Fatal(copy, err)
	}
	ticker, err := c.Markets.Ticker(ctx, "BTC-USDC")
	if err != nil || ticker.FundingRate.String() != "0.00001" {
		t.Fatal(ticker, err)
	}
	*ticker.Mark = decimal.MustParse("0")
	ticker, _ = c.Markets.Ticker(ctx, "BTC-USDC")
	if ticker.Mark.String() != "70000" {
		t.Fatal("shared pointer")
	}
	spot, _ := c.Markets.Ticker(ctx, "ETH-USDC")
	if spot.FundingRate != nil {
		t.Fatal("missing != zero")
	}
	trades, err := c.Markets.Trades(ctx, "BTC-USDC")
	if err != nil || trades[0].ID != "synthetic-1" {
		t.Fatal(trades, err)
	}
}
func TestDefaultPendingAndValidation(t *testing.T) {
	c, _ := txflow.NewClient()
	if _, err := c.Markets.List(context.Background()); !errors.Is(err, txflow.ErrPendingConfirmation) {
		t.Fatal(err)
	}
	demoClient, _ := demo.NewClient()
	if _, err := demoClient.Markets.Metadata(context.Background(), "UNKNOWN"); !errors.Is(err, txflow.ErrNotFound) {
		t.Fatal(err)
	}
	for _, req := range []types.OrderBookRequest{{Symbol: "", Depth: 1}, {Symbol: "BTC-USDC", Depth: -1}, {Symbol: "BTC-USDC", Depth: 1001}} {
		if _, err := demoClient.Markets.OrderBook(context.Background(), req); !errors.Is(err, txflow.ErrInvalidRequest) {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := demoClient.Markets.List(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := demoClient.Markets.List(nil); !errors.Is(err, txflow.ErrInvalidRequest) {
		t.Fatal(err)
	}
}
func TestConcurrentMarketReads(t *testing.T) {
	c, _ := demo.NewClient()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := c.Markets.OrderBook(context.Background(), types.OrderBookRequest{Symbol: "BTC-USDC"})
			if err != nil {
				t.Error(err)
			}
			_, err = c.Markets.Ticker(context.Background(), "BTC-USDC")
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

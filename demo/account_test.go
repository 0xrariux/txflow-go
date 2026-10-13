package demo_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	txflow "github.com/0xrariux/txflow-go"
	"github.com/0xrariux/txflow-go/decimal"
	"github.com/0xrariux/txflow-go/demo"
)

func TestAccountFixturesAndIsolation(t *testing.T) {
	c, _ := demo.NewClient()
	ctx := context.Background()
	state, err := c.Account.State(ctx, demo.Address)
	if err != nil || state.Equity.String() != "1010" {
		t.Fatal(state, err)
	}
	state.Balances[0].Asset = "changed"
	*state.Equity = decimal.MustParse("1")
	copy, _ := c.Account.State(ctx, demo.Address)
	if copy.Balances[0].Asset != "USDC" || copy.Equity.String() != "1010" {
		t.Fatal("shared state")
	}
	orders, err := c.Account.OpenOrders(ctx, demo.Address)
	if err != nil || len(orders) != 1 {
		t.Fatal(orders, err)
	}
	history, err := c.Account.HistoricalOrders(ctx, demo.Address)
	if err != nil || history[0].Status != "filled" {
		t.Fatal(history, err)
	}
	fills, err := c.Account.Fills(ctx, demo.Address)
	if err != nil || fills[0].Fee.String() != "0.1" {
		t.Fatal(fills, err)
	}
	funding, err := c.Account.FundingHistory(ctx, demo.Address)
	if err != nil || funding[0].Amount.String() != "-0.0001" {
		t.Fatal(funding, err)
	}
	if _, err := c.Account.State(ctx, "private-key"); !errors.Is(err, txflow.ErrInvalidRequest) {
		t.Fatal(err)
	}
	if _, err := c.Account.State(ctx, "0x1111111111111111111111111111111111111111"); !errors.Is(err, txflow.ErrNotFound) {
		t.Fatal(err)
	}
	live, _ := txflow.NewClient()
	if _, err := live.Account.State(ctx, demo.Address); !errors.Is(err, txflow.ErrPendingConfirmation) {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := c.Account.State(ctx, demo.Address)
			if err != nil {
				t.Error(err)
			}
			s.Balances[0].Asset = "caller-owned"
		}()
	}
	wg.Wait()
}

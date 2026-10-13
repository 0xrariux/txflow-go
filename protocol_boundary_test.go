package txflow

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/0xrariux/txflow-go/trading"
	"github.com/0xrariux/txflow-go/types"
)

type tripwireTransport struct{ calls atomic.Int32 }

func (t *tripwireTransport) RoundTrip(*http.Request) (*http.Response, error) {
	t.calls.Add(1)
	return nil, errors.New("unexpected I/O")
}
func TestUnconfirmedServicesNeverSendRequests(t *testing.T) {
	trap := &tripwireTransport{}
	c, err := NewClient(WithNetwork(Mainnet), WithBaseURL("http://127.0.0.1:1"), WithHTTPClient(&http.Client{Transport: trap}))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, e1 := c.Markets.List(ctx)
	_, e2 := c.Markets.OrderBook(ctx, types.OrderBookRequest{Symbol: "BTC-USDC"})
	_, e3 := c.Account.State(ctx, "0x0000000000000000000000000000000000000000")
	_, e4 := c.Stream.Connect(ctx)
	_, e5 := c.Trading.PlaceOrder(ctx, trading.OrderRequest{})
	e6 := c.Trading.CancelOrder(ctx, trading.CancelRequest{})
	e7 := c.Trading.CancelAll(ctx, trading.CancelAllRequest{})
	_, e8 := c.Trading.ModifyOrder(ctx, trading.ModifyRequest{})
	for _, err := range []error{e1, e2, e3, e4, e5, e6, e7, e8} {
		if !errors.Is(err, ErrPendingConfirmation) {
			t.Fatal(err)
		}
	}
	if trap.calls.Load() != 0 {
		t.Fatal("unsupported operation made HTTP requests")
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := c.Trading.PlaceOrder(canceled, trading.OrderRequest{}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

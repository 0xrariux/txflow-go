package demo

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xrariux/txflow-go/orderbook"
	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
	"io"
)

// RunBook demonstrates a synthetic snapshot/replay/recovery scenario.
// All sequence behavior here belongs to the local fixture, not TxFlow.
func RunBook(ctx context.Context, w io.Writer) error {
	c, err := NewClient()
	if err != nil {
		return err
	}
	snapshot, err := c.Markets.OrderBook(ctx, types.OrderBookRequest{Symbol: "BTC-USDC", Depth: 20})
	if err != nil {
		return err
	}
	book, err := orderbook.New(snapshot.Symbol, true)
	if err != nil {
		return err
	}
	if err := book.Load(snapshot); err != nil {
		return err
	}
	session, err := c.Stream.Connect(ctx)
	if err != nil {
		return err
	}
	defer session.Close()
	sub, err := session.SubscribeOrderBook(ctx, snapshot.Symbol)
	if err != nil {
		return err
	}
	defer sub.Close()
	if _, err := fmt.Fprintln(w, "SYNTHETIC OFFLINE REPLAY; contiguous sequences are a demo policy"); err != nil {
		return err
	}
	for event := range sub.Updates() {
		if event.Snapshot != nil {
			err = book.Load(*event.Snapshot)
		} else {
			err = book.Apply(*event.Delta)
		}
		if errors.Is(err, sdkerr.ErrStaleUpdate) {
			if _, err = fmt.Fprintln(w, "duplicate/stale delta ignored"); err != nil {
				return err
			}
			continue
		}
		if errors.Is(err, sdkerr.ErrSequenceGap) {
			if _, err = fmt.Fprintln(w, "sequence gap detected; book unavailable until recovery snapshot"); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		q, err := book.Quote()
		if err != nil {
			return err
		}
		top, err := book.Snapshot(2)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "sequence=%d bid=%s ask=%s mid=%s spread=%s top2=%v/%v\n", *top.Sequence, q.BestBid, q.BestAsk, q.Mid, q.Spread, top.Bids, top.Asks); err != nil {
			return err
		}
	}
	return sub.Err()
}

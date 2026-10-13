package stream

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

func events() []types.BookEvent {
	return []types.BookEvent{{Delta: &types.BookDelta{Symbol: "BTC", Sequence: 1}}, {Delta: &types.BookDelta{Symbol: "BTC", Sequence: 2}}}
}
func waitDone(t *testing.T, s *Subscription) {
	t.Helper()
	select {
	case <-s.Done():
	case <-time.After(time.Second):
		t.Fatal("replay did not stop")
	}
}
func TestReplayCopiesOrderAndCompletion(t *testing.T) {
	input := events()
	s, err := Replay(context.Background(), input, 0)
	if err != nil {
		t.Fatal(err)
	}
	input[0].Delta.Sequence = 99
	var seq []uint64
	for e := range s.Updates() {
		seq = append(seq, e.Delta.Sequence)
	}
	waitDone(t, s)
	if len(seq) != 2 || seq[0] != 1 || seq[1] != 2 || s.Err() != nil {
		t.Fatal(seq, s.Err())
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}
func TestBackpressureCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s, err := Replay(ctx, events(), 0)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	waitDone(t, s)
	if !errors.Is(s.Err(), context.Canceled) {
		t.Fatal(s.Err())
	}
	if _, open := <-s.Updates(); open {
		t.Fatal("channel not closed")
	}
	// Idempotent concurrent close must not panic or deadlock.
	s, _ = Replay(context.Background(), events(), 0)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _ = s.Close() }()
	}
	wg.Wait()
	waitDone(t, s)
}
func TestReplayInputValidation(t *testing.T) {
	for _, e := range [][]types.BookEvent{{{}}, {{Snapshot: &types.OrderBook{}, Delta: &types.BookDelta{}}}} {
		if _, err := Replay(context.Background(), e, 1); !errors.Is(err, sdkerr.ErrInvalidRequest) {
			t.Fatal(err)
		}
	}
	if _, err := Replay(nil, events(), 1); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal(err)
	}
	if _, err := Replay(context.Background(), events(), 4097); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal(err)
	}
	if _, err := New(nil).Connect(context.Background()); !errors.Is(err, sdkerr.ErrPendingConfirmation) {
		t.Fatal(err)
	}
}
func TestSessionShutdownAndParentCancellation(t *testing.T) {
	connector := &ReplayConnector{Symbol: "BTC", Events: events(), Buffer: 0}
	for _, closeSession := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		session, err := connector.Connect(ctx)
		if err != nil {
			t.Fatal(err)
		}
		sub, err := session.SubscribeOrderBook(context.Background(), "BTC")
		if err != nil {
			t.Fatal(err)
		}
		if closeSession {
			_ = session.Close()
		} else {
			cancel()
		}
		waitDone(t, sub)
		cancel()
		_ = session.Close()
		if _, err := session.SubscribeOrderBook(context.Background(), "BTC"); !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	}
}

func TestSubscriptionRegistryPrunesAndBounds(t *testing.T) {
	c := &ReplayConnector{Symbol: "BTC", Events: events(), Buffer: 0}
	raw, err := c.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	s := raw.(*replaySession)
	defer s.Close()
	for i := 0; i < 1024; i++ {
		if _, err := s.SubscribeOrderBook(context.Background(), "BTC"); err != nil {
			t.Fatal(i, err)
		}
	}
	if _, err := s.SubscribeOrderBook(context.Background(), "BTC"); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal("unbounded registry", err)
	}
	_ = s.subscriptions[0].Close()
	if _, err := s.SubscribeOrderBook(context.Background(), "BTC"); err != nil {
		t.Fatal("completed subscription not pruned", err)
	}
	if len(s.subscriptions) != 1024 {
		t.Fatal(len(s.subscriptions))
	}
}

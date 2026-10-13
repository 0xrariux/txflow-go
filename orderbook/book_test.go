package orderbook

import (
	"errors"
	"sync"
	"testing"

	"github.com/0xrariux/txflow-go/decimal"
	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

func level(p, s string) types.Level {
	return types.Level{Price: decimal.MustParse(p), Size: decimal.MustParse(s)}
}
func snapshot() types.OrderBook {
	seq := uint64(10)
	return types.OrderBook{Symbol: "BTC", Sequence: &seq, Bids: []types.Level{level("99", "2"), level("100", "1")}, Asks: []types.Level{level("102", "3"), level("101", "1")}}
}
func TestOrderingPrecisionAndCopies(t *testing.T) {
	b, _ := New("BTC", true)
	input := snapshot()
	if err := b.Load(input); err != nil {
		t.Fatal(err)
	}
	input.Bids[0] = level("1", "1")
	*input.Sequence = 0
	q, err := b.Quote()
	if err != nil || q.BestBid.String() != "100" || q.BestAsk.String() != "101" || q.Mid.String() != "100.5" || q.Spread.String() != "1" {
		t.Fatal(q, err)
	}
	out, err := b.Snapshot(1)
	if err != nil || len(out.Bids) != 1 || out.Bids[0].Price.String() != "100" {
		t.Fatal(out, err)
	}
	out.Bids[0] = level("1", "1")
	*out.Sequence = 0
	next, _ := b.Snapshot(1)
	if next.Bids[0].Price.String() != "100" || *next.Sequence != 10 {
		t.Fatal("shared book")
	}
	seq := uint64(1)
	tiny := types.OrderBook{Symbol: "BTC", Sequence: &seq, Bids: []types.Level{level("0.000000000000000001", "1")}, Asks: []types.Level{level("0.000000000000000002", "1")}}
	if err := b.Load(tiny); err != nil {
		t.Fatal(err)
	}
	q, _ = b.Quote()
	if q.Mid.String() != "0.0000000000000000015" {
		t.Fatal(q)
	}
}
func TestDeletionDuplicateGapAndRecovery(t *testing.T) {
	b, _ := New("BTC", true)
	if err := b.Load(snapshot()); err != nil {
		t.Fatal(err)
	}
	delta := types.BookDelta{Symbol: "BTC", Sequence: 11, Bids: []types.Level{level("100", "0")}, Asks: []types.Level{level("100.5", "2")}}
	if err := b.Apply(delta); err != nil {
		t.Fatal(err)
	}
	q, _ := b.Quote()
	if q.BestBid.String() != "99" || q.BestAsk.String() != "100.5" {
		t.Fatal(q)
	}
	if err := b.Apply(delta); !errors.Is(err, sdkerr.ErrStaleUpdate) {
		t.Fatal(err)
	}
	if err := b.Apply(types.BookDelta{Symbol: "BTC", Sequence: 13}); !errors.Is(err, sdkerr.ErrSequenceGap) {
		t.Fatal(err)
	}
	if _, err := b.Quote(); !errors.Is(err, sdkerr.ErrBookUnavailable) {
		t.Fatal(err)
	}
	if err := b.Apply(types.BookDelta{Symbol: "BTC", Sequence: 12}); !errors.Is(err, sdkerr.ErrBookUnavailable) {
		t.Fatal(err)
	}
	s := snapshot()
	seq := uint64(13)
	s.Sequence = &seq
	if err := b.Load(s); err != nil {
		t.Fatal(err)
	}
	if err := b.Apply(types.BookDelta{Symbol: "BTC", Sequence: 14}); err != nil {
		t.Fatal(err)
	}
}
func TestInvalidSnapshotsAndDeltas(t *testing.T) {
	b, _ := New("BTC", true)
	for _, s := range []types.OrderBook{{Symbol: "wrong"}, {Symbol: "BTC"}, {Symbol: "BTC", Sequence: new(uint64), Bids: []types.Level{level("1", "0")}}, {Symbol: "BTC", Sequence: new(uint64), Bids: []types.Level{level("1", "1"), level("1.0", "2")}}, {Symbol: "BTC", Sequence: new(uint64), Bids: []types.Level{level("2", "1")}, Asks: []types.Level{level("1", "1")}}} {
		if err := b.Load(s); !errors.Is(err, sdkerr.ErrInvalidRequest) {
			t.Fatal(err)
		}
	}
	if err := b.Load(snapshot()); err != nil {
		t.Fatal(err)
	}
	if err := b.Apply(types.BookDelta{Symbol: "BTC", Sequence: 11, Bids: []types.Level{level("101", "1")}}); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal(err)
	}
	if _, err := b.Snapshot(2); !errors.Is(err, sdkerr.ErrBookUnavailable) {
		t.Fatal("invalid delta kept book live", err)
	}
	if err := b.Load(snapshot()); err != nil {
		t.Fatal(err)
	}
	bad := snapshot()
	bad.Bids = make([]types.Level, MaxLevels+1)
	if err := b.Load(bad); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal(err)
	}
	if _, err := b.Snapshot(0); !errors.Is(err, sdkerr.ErrInvalidRequest) {
		t.Fatal(err)
	}
	empty := snapshot()
	empty.Asks = nil
	if err := b.Load(empty); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Quote(); !errors.Is(err, sdkerr.ErrBookUnavailable) {
		t.Fatal(err)
	}
}
func TestConcurrentBookReaders(t *testing.T) {
	b, _ := New("BTC", true)
	_ = b.Load(snapshot())
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_, err := b.Snapshot(2)
				if err != nil {
					t.Error(err)
				}
			}
		}()
	}
	for j := uint64(11); j < 31; j++ {
		if err := b.Apply(types.BookDelta{Symbol: "BTC", Sequence: j, Bids: []types.Level{level("100", "2")}}); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
}

func TestMalformedDeltaInvalidatesBook(t *testing.T) {
	for _, delta := range []types.BookDelta{
		{Symbol: "BTC", Sequence: 11, Bids: []types.Level{level("100", "-1")}},
		{Symbol: "BTC", Sequence: 11, Asks: []types.Level{level("101", "1"), level("101.0", "2")}},
	} {
		b, _ := New("BTC", true)
		_ = b.Load(snapshot())
		if err := b.Apply(delta); !errors.Is(err, sdkerr.ErrInvalidRequest) {
			t.Fatal(err)
		}
		if _, err := b.Quote(); !errors.Is(err, sdkerr.ErrBookUnavailable) {
			t.Fatal("malformed delta preserved usable quote", err)
		}
		if err := b.Load(snapshot()); err != nil {
			t.Fatal(err)
		}
		if _, err := b.Quote(); err != nil {
			t.Fatal(err)
		}
	}
}

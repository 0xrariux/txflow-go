// Package orderbook maintains a bounded local domain book. Its optional
// contiguous sequence policy is demonstrated on synthetic data only.
// It is not a verified TxFlow snapshot/delta implementation.
package orderbook

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/0xrariux/txflow-go/decimal"
	"github.com/0xrariux/txflow-go/markets"
	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

const MaxLevels = 10000

type Book struct {
	mu        sync.RWMutex
	symbol    string
	strict    bool
	valid     bool
	sequence  *uint64
	bids      map[string]types.Level
	asks      map[string]types.Level
	timestamp time.Time
}

func New(symbol string, strictSequence bool) (*Book, error) {
	if !markets.ValidSymbol(symbol) {
		return nil, sdkerr.ErrInvalidRequest
	}
	return &Book{symbol: symbol, strict: strictSequence}, nil
}
func levels(input []types.Level, allowDelete bool) (map[string]types.Level, error) {
	if len(input) > MaxLevels {
		return nil, sdkerr.ErrInvalidRequest
	}
	out := make(map[string]types.Level, len(input))
	for _, l := range input {
		key := l.Price.String()
		if l.Price.Sign() <= 0 || l.Size.Sign() < 0 || (!allowDelete && l.Size.Sign() == 0) {
			return nil, sdkerr.ErrInvalidRequest
		}
		if _, exists := out[key]; exists {
			return nil, fmt.Errorf("%w: duplicate price", sdkerr.ErrInvalidRequest)
		}
		out[key] = l
	}
	return out, nil
}
func uncrossed(bids, asks map[string]types.Level) bool {
	var bid, ask decimal.Decimal
	for _, l := range bids {
		if l.Price.Cmp(bid) > 0 {
			bid = l.Price
		}
	}
	for _, l := range asks {
		if ask.Sign() == 0 || l.Price.Cmp(ask) < 0 {
			ask = l.Price
		}
	}
	return bid.Sign() == 0 || ask.Sign() == 0 || bid.Cmp(ask) < 0
}
func (b *Book) Load(snapshot types.OrderBook) error {
	if snapshot.Symbol != b.symbol || (b.strict && snapshot.Sequence == nil) {
		return sdkerr.ErrInvalidRequest
	}
	bids, err := levels(snapshot.Bids, false)
	if err != nil {
		return err
	}
	asks, err := levels(snapshot.Asks, false)
	if err != nil {
		return err
	}
	if !uncrossed(bids, asks) {
		return fmt.Errorf("%w: crossed book", sdkerr.ErrInvalidRequest)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bids = bids
	b.asks = asks
	b.timestamp = snapshot.Timestamp
	b.valid = true
	b.sequence = nil
	if snapshot.Sequence != nil {
		s := *snapshot.Sequence
		b.sequence = &s
	}
	return nil
}
func copyLevels(src map[string]types.Level) map[string]types.Level {
	out := make(map[string]types.Level, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
func merge(dst, updates map[string]types.Level) {
	for k, v := range updates {
		if v.Size.Sign() == 0 {
			delete(dst, k)
		} else {
			dst[k] = v
		}
	}
}
func (b *Book) Apply(delta types.BookDelta) error {
	if delta.Symbol != b.symbol {
		return sdkerr.ErrInvalidRequest
	}
	bids, err := levels(delta.Bids, true)
	if err != nil {
		b.invalidate()
		return err
	}
	asks, err := levels(delta.Asks, true)
	if err != nil {
		b.invalidate()
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.valid {
		return sdkerr.ErrBookUnavailable
	}
	if b.strict {
		if delta.Sequence <= *b.sequence {
			return sdkerr.ErrStaleUpdate
		}
		if *b.sequence == ^uint64(0) || delta.Sequence != *b.sequence+1 {
			b.valid = false
			return sdkerr.ErrSequenceGap
		}
	}
	nextBids, nextAsks := copyLevels(b.bids), copyLevels(b.asks)
	merge(nextBids, bids)
	merge(nextAsks, asks)
	if len(nextBids) > MaxLevels || len(nextAsks) > MaxLevels || !uncrossed(nextBids, nextAsks) {
		// An unusable delta cannot be skipped while pretending continuity.
		b.valid = false
		return fmt.Errorf("%w: invalid delta; snapshot required", sdkerr.ErrInvalidRequest)
	}
	b.bids = nextBids
	b.asks = nextAsks
	b.timestamp = delta.Timestamp
	s := delta.Sequence
	b.sequence = &s
	return nil
}
func sorted(src map[string]types.Level, descending bool, depth int) []types.Level {
	out := make([]types.Level, 0, len(src))
	for _, v := range src {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		n := out[i].Price.Cmp(out[j].Price)
		if descending {
			return n > 0
		}
		return n < 0
	})
	if len(out) > depth {
		out = out[:depth]
	}
	return out
}
func (b *Book) Snapshot(depth int) (types.OrderBook, error) {
	if depth < 1 || depth > MaxLevels {
		return types.OrderBook{}, sdkerr.ErrInvalidRequest
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if !b.valid {
		return types.OrderBook{}, sdkerr.ErrBookUnavailable
	}
	out := types.OrderBook{Symbol: b.symbol, Bids: sorted(b.bids, true, depth), Asks: sorted(b.asks, false, depth), Timestamp: b.timestamp}
	if b.sequence != nil {
		s := *b.sequence
		out.Sequence = &s
	}
	return out, nil
}

type Quote struct{ BestBid, BestAsk, Mid, Spread decimal.Decimal }

func (b *Book) Quote() (Quote, error) {
	book, err := b.Snapshot(1)
	if err != nil {
		return Quote{}, err
	}
	if len(book.Bids) == 0 || len(book.Asks) == 0 {
		return Quote{}, sdkerr.ErrBookUnavailable
	}
	bid, ask := book.Bids[0].Price, book.Asks[0].Price
	return Quote{bid, ask, bid.Add(ask).Half(), ask.Sub(bid)}, nil
}

func (b *Book) invalidate() { b.mu.Lock(); b.valid = false; b.mu.Unlock() }

package types

import "time"

// BookDelta is a proposed domain delta, used only by synthetic replay today.
// A zero size removes a price level. Sequence semantics are adapter-defined.
type BookDelta struct {
	Symbol    string    `json:"symbol"`
	Sequence  uint64    `json:"sequence"`
	Bids      []Level   `json:"bids"`
	Asks      []Level   `json:"asks"`
	Timestamp time.Time `json:"timestamp"`
}
type BookEvent struct {
	Snapshot *OrderBook `json:"snapshot,omitempty"`
	Delta    *BookDelta `json:"delta,omitempty"`
}

func CloneBook(b OrderBook) OrderBook {
	b.Bids = append([]Level(nil), b.Bids...)
	b.Asks = append([]Level(nil), b.Asks...)
	if b.Sequence != nil {
		seq := *b.Sequence
		b.Sequence = &seq
	}
	return b
}
func CloneBookEvent(e BookEvent) BookEvent {
	if e.Snapshot != nil {
		v := CloneBook(*e.Snapshot)
		e.Snapshot = &v
	}
	if e.Delta != nil {
		v := *e.Delta
		v.Bids = append([]Level(nil), v.Bids...)
		v.Asks = append([]Level(nil), v.Asks...)
		e.Delta = &v
	}
	return e
}

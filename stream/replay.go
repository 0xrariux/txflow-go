package stream

import (
	"context"
	"sync"

	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

// Subscription owns one producer. Backpressure blocks the producer, bounded
// by context cancellation; no updates are silently dropped.
type Subscription struct {
	updates chan types.BookEvent
	done    chan struct{}
	cancel  context.CancelFunc
	mu      sync.Mutex
	err     error
}

func (s *Subscription) Updates() <-chan types.BookEvent { return s.updates }
func (s *Subscription) Done() <-chan struct{}           { return s.done }
func (s *Subscription) Err() error                      { s.mu.Lock(); defer s.mu.Unlock(); return s.err }
func (s *Subscription) Close() error                    { s.cancel(); <-s.done; return nil }

// Replay copies events before starting; callers must not mutate inputs during
// this call. Its JSON/event format is local domain data, not a WS wire schema.
func Replay(ctx context.Context, events []types.BookEvent, buffer int) (*Subscription, error) {
	if ctx == nil || buffer < 0 || buffer > 4096 || len(events) > 65536 {
		return nil, sdkerr.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	input := make([]types.BookEvent, len(events))
	for i, e := range events {
		if e.Snapshot != nil && (len(e.Snapshot.Bids) > 10000 || len(e.Snapshot.Asks) > 10000) {
			return nil, sdkerr.ErrInvalidRequest
		}
		if e.Delta != nil && (len(e.Delta.Bids) > 10000 || len(e.Delta.Asks) > 10000) {
			return nil, sdkerr.ErrInvalidRequest
		}
		if (e.Snapshot == nil) == (e.Delta == nil) {
			return nil, sdkerr.ErrInvalidRequest
		}
		input[i] = types.CloneBookEvent(e)
	}
	ctx, cancel := context.WithCancel(ctx)
	s := &Subscription{updates: make(chan types.BookEvent, buffer), done: make(chan struct{}), cancel: cancel}
	go func() {
		defer close(s.done)
		defer close(s.updates)
		defer cancel()
		setErr := func(err error) { s.mu.Lock(); s.err = err; s.mu.Unlock() }
		for _, event := range input {
			if err := ctx.Err(); err != nil {
				setErr(err)
				return
			}
			select {
			case <-ctx.Done():
				setErr(ctx.Err())
				return
			case s.updates <- event:
			}
		}
	}()
	return s, nil
}

// ReplayConnector creates finite local sessions for synthetic demonstrations.
// Its configuration must remain immutable after Connect. Sessions bound the
// active subscription registry to 1024 entries and prune completed producers.
type ReplayConnector struct {
	Symbol string
	Events []types.BookEvent
	Buffer int
}

func (c *ReplayConnector) Connect(ctx context.Context) (Session, error) {
	if ctx == nil || ctx.Err() != nil {
		if ctx != nil {
			return nil, ctx.Err()
		}
		return nil, sdkerr.ErrInvalidRequest
	}
	ctx, cancel := context.WithCancel(ctx)
	return &replaySession{ctx: ctx, cancel: cancel, connector: c}, nil
}

type replaySession struct {
	ctx           context.Context
	cancel        context.CancelFunc
	connector     *ReplayConnector
	mu            sync.Mutex
	closed        bool
	subscriptions []*Subscription
}

func (s *replaySession) SubscribeOrderBook(ctx context.Context, symbol string) (*Subscription, error) {
	if ctx == nil || symbol != s.connector.Symbol {
		return nil, sdkerr.ErrInvalidRequest
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, context.Canceled
	}
	if err := s.ctx.Err(); err != nil {
		return nil, err
	}
	active := s.subscriptions[:0]
	for _, sub := range s.subscriptions {
		select {
		case <-sub.Done():
		default:
			active = append(active, sub)
		}
	}
	for i := len(active); i < len(s.subscriptions); i++ {
		s.subscriptions[i] = nil
	}
	s.subscriptions = active
	if len(active) >= 1024 {
		return nil, sdkerr.ErrInvalidRequest
	}
	child, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(s.ctx, cancel)
	sub, err := Replay(child, s.connector.Events, s.connector.Buffer)
	if err != nil {
		stop()
		cancel()
		return nil, err
	}
	s.subscriptions = append(s.subscriptions, sub)
	go func() { <-sub.Done(); stop(); cancel() }()
	return sub, nil
}
func (s *replaySession) Close() error {
	s.mu.Lock()
	s.closed = true
	s.cancel()
	subs := append([]*Subscription(nil), s.subscriptions...)
	s.mu.Unlock()
	for _, sub := range subs {
		<-sub.Done()
	}
	return nil
}

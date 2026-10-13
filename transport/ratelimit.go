package transport

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xrariux/txflow-go/sdkerr"
)

// IntervalLimiter spaces attempts across callers. It is not a TxFlow quota.
// Canceled reservations may leave a conservative idle slot.
type IntervalLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

func NewIntervalLimiter(interval time.Duration) (*IntervalLimiter, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("%w: positive interval required", sdkerr.ErrInvalidRequest)
	}
	return &IntervalLimiter{interval: interval}, nil
}
func (l *IntervalLimiter) Wait(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	now := time.Now()
	when := l.next
	if when.Before(now) {
		when = now
	}
	l.next = when.Add(l.interval)
	l.mu.Unlock()
	return wait(ctx, time.Until(when))
}

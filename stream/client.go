// Package stream defines a future streaming boundary and local replay tools.
// No live TxFlow WebSocket URL, handshake, heartbeat or recovery is assumed.
package stream

import (
	"context"
	"github.com/0xrariux/txflow-go/sdkerr"
)

type Connector interface {
	Connect(context.Context) (Session, error)
}
type Session interface {
	SubscribeOrderBook(context.Context, string) (*Subscription, error)
	Close() error
}
type Client struct{ connector Connector }

func New(c Connector) *Client { return &Client{connector: c} }
func (c *Client) Connect(ctx context.Context) (Session, error) {
	if ctx == nil {
		return nil, sdkerr.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c.connector == nil {
		return nil, sdkerr.ErrPendingConfirmation
	}
	return c.connector.Connect(ctx)
}

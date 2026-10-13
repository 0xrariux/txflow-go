// Package txflow is a community SDK prototype. Live protocol integration is
// pending TxFlow confirmation. Use package demo explicitly for synthetic data.
package txflow

import (
	"fmt"

	"github.com/0xrariux/txflow-go/account"
	"github.com/0xrariux/txflow-go/markets"
	"github.com/0xrariux/txflow-go/stream"
	"github.com/0xrariux/txflow-go/trading"
	"github.com/0xrariux/txflow-go/transport"
)

type Client struct {
	Trading *trading.Client
	Stream  *stream.Client
	Markets *markets.Client
	Account *account.Client
	// Transport is generic HTTP infrastructure, present only with WithBaseURL.
	// It has no TxFlow endpoint mapping.
	Transport *transport.Client
	cfg       config
}

func NewClient(opts ...Option) (*Client, error) {
	cfg := config{}
	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("%w: nil option", ErrInvalidRequest)
		}
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	c := &Client{
		cfg:     cfg,
		Markets: markets.New(cfg.markets),
		Account: account.New(cfg.account),
		Stream:  stream.New(cfg.stream),
		Trading: &trading.Client{},
	}
	if cfg.transport.BaseURL != "" {
		t, err := transport.New(cfg.transport)
		if err != nil {
			return nil, err
		}
		c.Transport = t
	}
	return c, nil
}
func (c *Client) Network() Network { return c.cfg.network }

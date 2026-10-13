// Package trading defines a proposed API only. Every operation fails locally
// pending TxFlow confirmation. It contains no exchange endpoint or signer.
package trading

import (
	"context"
	"github.com/0xrariux/txflow-go/decimal"
	"github.com/0xrariux/txflow-go/sdkerr"
	"github.com/0xrariux/txflow-go/types"
)

type OrderType string

const (
	Limit  OrderType = "limit"
	Market OrderType = "market"
)

type TimeInForce string

const (
	GoodTillCanceled  TimeInForce = "gtc"
	ImmediateOrCancel TimeInForce = "ioc"
	PostOnly          TimeInForce = "post_only"
)

type OrderRequest struct {
	Symbol        string
	Side          types.Side
	Type          OrderType
	TimeInForce   TimeInForce
	ReduceOnly    bool
	Price         *decimal.Decimal
	Quantity      decimal.Decimal
	ClientOrderID string
}
type OrderResponse struct {
	OrderID       string
	ClientOrderID string
	Status        string
}
type CancelRequest struct {
	Symbol        string
	OrderID       string
	ClientOrderID string
}
type CancelAllRequest struct{ Symbol string }
type ModifyRequest struct {
	OrderID     string
	Replacement OrderRequest
}
type Service interface {
	PlaceOrder(context.Context, OrderRequest) (*OrderResponse, error)
	CancelOrder(context.Context, CancelRequest) error
	CancelAll(context.Context, CancelAllRequest) error
	ModifyOrder(context.Context, ModifyRequest) (*OrderResponse, error)
}

// Client is deliberately not configurable with a signer or submitter.
type Client struct{}

func pending(ctx context.Context) error {
	if ctx == nil {
		return sdkerr.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return sdkerr.ErrPendingConfirmation
}
func (*Client) PlaceOrder(ctx context.Context, _ OrderRequest) (*OrderResponse, error) {
	return nil, pending(ctx)
}
func (*Client) CancelOrder(ctx context.Context, _ CancelRequest) error  { return pending(ctx) }
func (*Client) CancelAll(ctx context.Context, _ CancelAllRequest) error { return pending(ctx) }
func (*Client) ModifyOrder(ctx context.Context, _ ModifyRequest) (*OrderResponse, error) {
	return nil, pending(ctx)
}

var _ Service = (*Client)(nil)

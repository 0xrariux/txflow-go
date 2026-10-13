package txflow

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/0xrariux/txflow-go/account"
	"github.com/0xrariux/txflow-go/markets"
	"github.com/0xrariux/txflow-go/stream"
	"github.com/0xrariux/txflow-go/transport"
)

type Network string

const (
	Mainnet Network = "mainnet"
	Testnet Network = "testnet"
)

type config struct {
	stream       stream.Connector
	account      account.Reader
	markets      markets.Reader
	network      Network
	websocketURL string
	transport    transport.Config
}
type Option func(*config) error

func WithNetwork(n Network) Option {
	return func(c *config) error {
		if n != Mainnet && n != Testnet {
			return fmt.Errorf("%w: unknown network", ErrInvalidRequest)
		}
		c.network = n
		return nil
	}
}
func WithBaseURL(raw string) Option {
	return func(c *config) error {
		if _, err := transport.ValidateURL(raw); err != nil {
			return err
		}
		c.transport.BaseURL = raw
		return nil
	}
}
func WithWebSocketURL(raw string) Option {
	return func(c *config) error {
		u, err := url.Parse(raw)
		if err != nil || u.User != nil || u.Fragment != "" || u.RawQuery != "" {
			return fmt.Errorf("%w: invalid WebSocket URL", ErrInvalidRequest)
		}
		translated := strings.Replace(strings.Replace(raw, "wss://", "https://", 1), "ws://", "http://", 1)
		if u.Scheme != "ws" && u.Scheme != "wss" {
			return fmt.Errorf("%w: invalid WebSocket scheme", ErrInvalidRequest)
		}
		if _, err := transport.ValidateURL(translated); err != nil {
			return err
		}
		c.websocketURL = raw
		return nil
	}
}
func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return fmt.Errorf("%w: positive timeout required", ErrInvalidRequest)
		}
		c.transport.Timeout = d
		return nil
	}
}
func WithHTTPClient(h *http.Client) Option {
	return func(c *config) error {
		if h == nil {
			return fmt.Errorf("%w: nil HTTP client", ErrInvalidRequest)
		}
		c.transport.HTTPClient = h
		return nil
	}
}
func WithUserAgent(s string) Option {
	return func(c *config) error {
		if s == "" || strings.ContainsAny(s, "\r\n") {
			return fmt.Errorf("%w: invalid user-agent", ErrInvalidRequest)
		}
		c.transport.UserAgent = s
		return nil
	}
}
func WithRetries(n int) Option {
	return func(c *config) error {
		if n < 0 || n > 5 {
			return fmt.Errorf("%w: retries must be 0..5", ErrInvalidRequest)
		}
		c.transport.MaxRetries = n
		return nil
	}
}
func WithLimiter(l transport.Limiter) Option {
	return func(c *config) error {
		if nilValue(l) {
			return fmt.Errorf("%w: nil limiter", ErrInvalidRequest)
		}
		c.transport.Limiter = l
		return nil
	}
}
func WithObserver(h func(transport.Event)) Option {
	return func(c *config) error { c.transport.Observe = h; return nil }
}

// WithMarketReader explicitly injects an authorized or offline domain adapter.
func WithMarketReader(r markets.Reader) Option {
	return func(c *config) error {
		if nilValue(r) {
			return fmt.Errorf("%w: nil market reader", ErrInvalidRequest)
		}
		c.markets = r
		return nil
	}
}

// WithAccountReader explicitly injects an authorized or offline query adapter.
func WithAccountReader(r account.Reader) Option {
	return func(c *config) error {
		if nilValue(r) {
			return fmt.Errorf("%w: nil account reader", ErrInvalidRequest)
		}
		c.account = r
		return nil
	}
}

// WithStreamConnector injects an explicit local or authorized connector.
func WithStreamConnector(s stream.Connector) Option {
	return func(c *config) error {
		if nilValue(s) {
			return fmt.Errorf("%w: nil stream connector", ErrInvalidRequest)
		}
		c.stream = s
		return nil
	}
}

func nilValue(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Func, reflect.Map, reflect.Slice, reflect.Chan:
		return v.IsNil()
	}
	return false
}

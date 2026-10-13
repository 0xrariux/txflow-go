package txflow

import (
	"errors"
	"github.com/0xrariux/txflow-go/account"
	"github.com/0xrariux/txflow-go/markets"
	"github.com/0xrariux/txflow-go/stream"
	"github.com/0xrariux/txflow-go/transport"
	"net/http"
	"testing"
	"time"
)

func TestConfiguration(t *testing.T) {
	c, err := NewClient()
	if err != nil || c.Transport != nil || c.Network() != "" {
		t.Fatal(c, err)
	}
	for _, opt := range []Option{nil, WithNetwork("invalid"), WithTimeout(0), WithHTTPClient(nil), WithRetries(6), WithUserAgent("a\nb"), WithBaseURL("http://example.com"), WithWebSocketURL("https://example.com")} {
		if _, err := NewClient(opt); !errors.Is(err, ErrInvalidRequest) {
			t.Fatal(err)
		}
	}
	c, err = NewClient(WithNetwork(Testnet), WithBaseURL("http://127.0.0.1:1234"), WithWebSocketURL("ws://localhost:1234"), WithTimeout(time.Second), WithHTTPClient(&http.Client{}), WithRetries(2))
	if err != nil || c.Network() != Testnet || c.Transport == nil {
		t.Fatal(c, err)
	}
}

func TestTypedNilAdaptersRejected(t *testing.T) {
	var marketReader *markets.Client
	var accountReader *account.Client
	var connector *stream.ReplayConnector
	var limiter *transport.IntervalLimiter
	for _, opt := range []Option{WithMarketReader(marketReader), WithAccountReader(accountReader), WithStreamConnector(connector), WithLimiter(limiter)} {
		if _, err := NewClient(opt); !errors.Is(err, ErrInvalidRequest) {
			t.Fatal("typed nil accepted", err)
		}
	}
}

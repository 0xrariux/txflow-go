package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xrariux/txflow-go/sdkerr"
)

func newTestClient(t *testing.T, h http.HandlerFunc, adjust func(*Config)) (*Client, *httptest.Server) {
	t.Helper()
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)
	cfg := Config{BaseURL: s.URL, Timeout: time.Second, InitialBackoff: time.Millisecond, MaxBackoff: 2 * time.Millisecond}
	if adjust != nil {
		adjust(&cfg)
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return c, s
}
func TestJSONReadAndHeaders(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/read" || r.URL.Query().Get("a") != "b" || r.Header.Get("User-Agent") != "txflow-go/0.1.0-prototype" || r.Header.Get("Accept") != "application/json" {
			t.Error("request metadata missing")
		}
		fmt.Fprint(w, `{"quantity":"0.000000000000000001","extra":true}`)
	}, nil)
	var out struct{ Quantity string }
	if err := c.DoJSON(context.Background(), Request{Path: "/read?a=b"}, &out); err != nil {
		t.Fatal(err)
	}
	if out.Quantity != "0.000000000000000001" {
		t.Fatal(out)
	}
}
func TestErrorsAndMalformedBodies(t *testing.T) {
	cases := []struct {
		name, body string
		status     int
		want       error
	}{
		{"rate", "secret", 429, sdkerr.ErrRateLimited},
		{"auth", "secret", 403, sdkerr.ErrUnauthorized},
		{"missing", "", 404, sdkerr.ErrNotFound},
		{"invalid", "", 422, sdkerr.ErrInvalidRequest},
		{"broken", `{"x":`, 200, sdkerr.ErrMalformedResponse},
		{"trailing", `{} {}`, 200, sdkerr.ErrMalformedResponse},
		{"garbage", `{} garbage`, 200, sdkerr.ErrMalformedResponse},
		{"empty", "", 200, sdkerr.ErrMalformedResponse},
		{"null", "null", 200, sdkerr.ErrMalformedResponse},
		{"oversized", strings.Repeat("x", 33), 200, sdkerr.ErrResponseTooLarge},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Request-ID", "test-id")
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			}, func(cfg *Config) { cfg.MaxResponseBytes = 32 })
			var out map[string]any
			err := c.DoJSON(context.Background(), Request{Path: "/"}, &out)
			if !errors.Is(err, tc.want) {
				t.Fatalf("got %v want %v", err, tc.want)
			}
			if tc.status != 200 {
				var api *sdkerr.APIError
				if !errors.As(err, &api) || api.StatusCode != tc.status || api.RequestID != "test-id" {
					t.Fatal(err)
				}
				if strings.Contains(err.Error(), "secret") {
					t.Fatal("server body exposed")
				}
			}
		})
	}
}
func TestReadRetriesAndNoWriteRetries(t *testing.T) {
	for _, tc := range []struct {
		method   string
		read     bool
		expected int
	}{{"GET", false, 3}, {"POST", true, 3}, {"POST", false, 1}} {
		t.Run(fmt.Sprintf("%s-%t", tc.method, tc.read), func(t *testing.T) {
			var calls atomic.Int32
			c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				n := calls.Add(1)
				if tc.method == "POST" && r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing JSON content type")
				}
				if n < 3 {
					w.WriteHeader(503)
					return
				}
				fmt.Fprint(w, `{"ok":true}`)
			}, func(cfg *Config) { cfg.MaxRetries = 2 })
			var out map[string]any
			req := Request{Method: tc.method, Path: "/", ReadOnly: tc.read}
			if tc.method == "POST" {
				req.Body = map[string]string{"query": "synthetic"}
			}
			err := c.DoJSON(context.Background(), req, &out)
			if int(calls.Load()) != tc.expected {
				t.Fatal(calls.Load())
			}
			if (err == nil) != (tc.expected == 3) {
				t.Fatal(err)
			}
		})
	}
}
func TestRetryBoundAndPermanentFailures(t *testing.T) {
	for _, status := range []int{429, 503, 400, 500} {
		var calls atomic.Int32
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(status) }, func(cfg *Config) { cfg.MaxRetries = 2 })
		err := c.DoJSON(context.Background(), Request{Path: "/"}, nil)
		want := int32(1)
		if status == 429 || status == 503 {
			want = 3
		}
		if err == nil || calls.Load() != want {
			t.Fatalf("status %d attempts %d err %v", status, calls.Load(), err)
		}
	}
}
func TestRetryAfterDoesNotRetryEarly(t *testing.T) {
	var calls atomic.Int32
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(429)
	}, func(cfg *Config) { cfg.MaxRetries = 2; cfg.Timeout = 30 * time.Millisecond })
	err := c.DoJSON(context.Background(), Request{Path: "/"}, nil)
	if !errors.Is(err, sdkerr.ErrRateLimited) || calls.Load() != 1 {
		t.Fatal(err, calls.Load())
	}
	if d, ok := retryAfter("1", time.Now()); !ok || d != time.Second {
		t.Fatal(d)
	}
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if d, ok := retryAfter(now.Add(time.Minute).Format(http.TimeFormat), now); !ok || d != time.Minute {
		t.Fatal(d)
	}
}
func TestCancellationAndWholeOperationTimeout(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }, func(cfg *Config) { cfg.Timeout = 20 * time.Millisecond })
	err := c.DoJSON(context.Background(), Request{Path: "/"}, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := c.DoJSON(ctx, Request{Path: "/"}, nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	var calls atomic.Int32
	c, _ = newTestClient(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(503) }, func(cfg *Config) {
		cfg.MaxRetries = 5
		cfg.Timeout = 10 * time.Millisecond
		cfg.InitialBackoff = time.Second
		cfg.MaxBackoff = time.Second
	})
	if err := c.DoJSON(context.Background(), Request{Path: "/"}, nil); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	if calls.Load() != 1 {
		t.Fatal("deadline did not bound backoff")
	}
}
func TestRedirectAndOriginProtection(t *testing.T) {
	var targetCalls atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { targetCalls.Add(1) }))
	defer target.Close()
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, 302) }, nil)
	err := c.DoJSON(context.Background(), Request{Path: "/"}, nil)
	var api *sdkerr.APIError
	if !errors.As(err, &api) || api.StatusCode != 302 || targetCalls.Load() != 0 {
		t.Fatal(err)
	}
	for _, p := range []string{target.URL, "//other.test/", "/%2e%2e/", "/#secret"} {
		if err := c.DoJSON(context.Background(), Request{Path: p}, nil); !errors.Is(err, sdkerr.ErrInvalidRequest) {
			t.Fatal(p, err)
		}
	}
}
func TestConfigAndDestinationValidation(t *testing.T) {
	for _, s := range []string{"", "ftp://a", "http://example.com", "https://user:pass@example.com", "https://a/#x", "https://a/?key=secret"} {
		if _, err := New(Config{BaseURL: s}); !errors.Is(err, sdkerr.ErrInvalidRequest) {
			t.Fatal(s, err)
		}
	}
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Error("invalid input reached server") }, nil)
	var nilDst *map[string]any
	for _, dst := range []any{map[string]any{}, nilDst} {
		if err := c.DoJSON(context.Background(), Request{Path: "/"}, dst); !errors.Is(err, sdkerr.ErrInvalidRequest) {
			t.Fatal(err)
		}
	}
}
func TestObserverConcurrencyAndCustomClientNotMutated(t *testing.T) {
	var mu sync.Mutex
	var events []Event
	custom := &http.Client{Timeout: time.Minute}
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `{}`) }, func(cfg *Config) {
		cfg.HTTPClient = custom
		cfg.Observe = func(e Event) { mu.Lock(); events = append(events, e); mu.Unlock() }
	})
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var out map[string]any
			if err := c.DoJSON(context.Background(), Request{Path: "/?private=secret"}, &out); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if custom.CheckRedirect != nil || custom.Timeout != time.Minute {
		t.Fatal("custom client mutated")
	}
	if len(events) != 20 || !reflect.DeepEqual(events[0].Method, "GET") {
		t.Fatal(events)
	}
}
func TestLimiterCancellation(t *testing.T) {
	l, err := NewIntervalLimiter(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := l.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
}

func TestOverflowRetryAfterIsConservative(t *testing.T) {
	d, ok := retryAfter("999999999999999999999999999999", time.Now())
	if !ok || d != time.Duration(1<<63-1) {
		t.Fatal(d, ok)
	}
}

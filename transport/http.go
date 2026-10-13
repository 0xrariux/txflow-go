// Package transport supplies generic HTTP infrastructure. It defines no
// TxFlow endpoints, wire envelopes, authentication or numeric rate quotas.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/0xrariux/txflow-go/sdkerr"
)

type Limiter interface{ Wait(context.Context) error }

// Event contains no payloads, URLs, addresses, request IDs or secrets.
// Hooks run synchronously and must be fast and concurrency-safe.
type Event struct {
	Method     string
	Attempt    int
	StatusCode int
	Duration   time.Duration
	Outcome    string
}
type Config struct {
	BaseURL          string
	HTTPClient       *http.Client
	Timeout          time.Duration
	UserAgent        string
	MaxRetries       int
	InitialBackoff   time.Duration
	MaxBackoff       time.Duration
	MaxResponseBytes int64
	Limiter          Limiter
	Observe          func(Event)
}
type Client struct {
	base *url.URL
	http *http.Client
	cfg  Config
}

// Request is generic adapter infrastructure. ReadOnly permits retries for
// a verified read-only POST. It must never be used for trading submissions.
type Request struct {
	Method   string
	Path     string
	Body     any
	ReadOnly bool
}

func ValidateURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil || u == nil || u.Host == "" || u.User != nil || u.Fragment != "" || u.RawQuery != "" || u.Opaque != "" {
		return nil, fmt.Errorf("%w: invalid base URL", sdkerr.ErrInvalidRequest)
	}
	if u.Scheme != "https" {
		ip := net.ParseIP(u.Hostname())
		if u.Scheme != "http" || (u.Hostname() != "localhost" && (ip == nil || !ip.IsLoopback())) {
			return nil, fmt.Errorf("%w: HTTPS required outside loopback", sdkerr.ErrInvalidRequest)
		}
	}
	return u, nil
}
func New(cfg Config) (*Client, error) {
	base, err := ValidateURL(cfg.BaseURL)
	if err != nil {
		return nil, err
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.InitialBackoff == 0 {
		cfg.InitialBackoff = 100 * time.Millisecond
	}
	if cfg.MaxBackoff == 0 {
		cfg.MaxBackoff = 2 * time.Second
	}
	if cfg.MaxResponseBytes == 0 {
		cfg.MaxResponseBytes = 4 << 20
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "txflow-go/0.1.0-prototype"
	}
	if cfg.Timeout < 0 || cfg.MaxRetries < 0 || cfg.MaxRetries > 5 || cfg.InitialBackoff < 0 || cfg.MaxBackoff < cfg.InitialBackoff || cfg.MaxBackoff > time.Minute || cfg.MaxResponseBytes < 1 || cfg.MaxResponseBytes > 64<<20 || strings.ContainsAny(cfg.UserAgent, "\r\n") {
		return nil, fmt.Errorf("%w: invalid transport configuration", sdkerr.ErrInvalidRequest)
	}
	h := http.Client{}
	if cfg.HTTPClient != nil {
		h = *cfg.HTTPClient
	}
	h.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	// The per-operation context also covers throttling and backoff.
	return &Client{base: base, http: &h, cfg: cfg}, nil
}
func (c *Client) DoJSON(ctx context.Context, r Request, dst any) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", sdkerr.ErrInvalidRequest)
	}
	if r.Method == "" {
		r.Method = http.MethodGet
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		return fmt.Errorf("%w: only GET/POST supported", sdkerr.ErrInvalidRequest)
	}
	if r.Method == http.MethodGet && r.Body != nil {
		return fmt.Errorf("%w: GET body", sdkerr.ErrInvalidRequest)
	}
	v := reflect.ValueOf(dst)
	if dst != nil && (v.Kind() != reflect.Pointer || v.IsNil()) {
		return fmt.Errorf("%w: response must be non-nil pointer", sdkerr.ErrInvalidRequest)
	}
	ref, err := url.Parse(r.Path)
	if err != nil || ref == nil || ref.IsAbs() || ref.Host != "" || ref.User != nil || ref.Fragment != "" || strings.Contains(ref.Path, "..") {
		return fmt.Errorf("%w: path must stay on configured origin", sdkerr.ErrInvalidRequest)
	}
	endpoint := c.base.ResolveReference(ref)
	var body []byte
	if r.Body != nil {
		body, err = json.Marshal(r.Body)
		if err != nil {
			return fmt.Errorf("%w: encode JSON", sdkerr.ErrInvalidRequest)
		}
	}
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()
	safe := r.Method == http.MethodGet || r.ReadOnly
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if c.cfg.Limiter != nil {
			if err := c.cfg.Limiter.Wait(ctx); err != nil {
				return err
			}
		}
		req, err := http.NewRequestWithContext(ctx, r.Method, endpoint.String(), bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("%w: construct HTTP request", sdkerr.ErrInvalidRequest)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.cfg.UserAgent)
		if r.Body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		started := time.Now()
		resp, err := c.http.Do(req)
		status := 0
		outcome := "network_error"
		var advice string
		if err == nil {
			status = resp.StatusCode
			advice = resp.Header.Get("Retry-After")
			err = c.decode(resp, dst)
			outcome = "ok"
			if err != nil {
				outcome = "response_error"
			}
		}
		if c.cfg.Observe != nil {
			c.cfg.Observe(Event{r.Method, attempt + 1, status, time.Since(started), outcome})
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			return nil
		}
		retry := status == 429 || status == 502 || status == 503 || status == 504
		if status == 0 {
			var ne net.Error
			retry = errors.As(err, &ne)
		}
		if !safe || !retry || attempt >= c.cfg.MaxRetries {
			return err
		}
		delay := c.cfg.InitialBackoff
		for i := 0; i < attempt && delay < c.cfg.MaxBackoff; i++ {
			if delay > c.cfg.MaxBackoff/2 {
				delay = c.cfg.MaxBackoff
			} else {
				delay *= 2
			}
		}
		// Downward jitter stays within the configured cap.
		delay = delay/2 + time.Duration(rand.Int63n(int64(delay-delay/2)+1))
		if advice != "" {
			advised, valid := retryAfter(advice, time.Now())
			if valid && advised > delay {
				if deadline, ok := ctx.Deadline(); ok && advised >= time.Until(deadline) {
					return err
				}
				delay = advised
			}
		}
		if err := wait(ctx, delay); err != nil {
			return err
		}
	}
}
func (c *Client) decode(resp *http.Response, dst any) error {
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Bounded drain for connection reuse; never expose arbitrary response bodies.
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
		id := resp.Header.Get("X-Request-ID")
		if len(id) > 256 {
			id = ""
		}
		return &sdkerr.APIError{StatusCode: resp.StatusCode, RequestID: id}
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, c.cfg.MaxResponseBytes+1))
	if err != nil {
		return fmt.Errorf("%w: read body", sdkerr.ErrMalformedResponse)
	}
	if int64(len(raw)) > c.cfg.MaxResponseBytes {
		return sdkerr.ErrResponseTooLarge
	}
	if dst == nil {
		return nil
	}
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return fmt.Errorf("%w: null response", sdkerr.ErrMalformedResponse)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%w: invalid JSON or model", sdkerr.ErrMalformedResponse)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return fmt.Errorf("%w: trailing JSON", sdkerr.ErrMalformedResponse)
	}
	return nil
}
func retryAfter(s string, now time.Time) (time.Duration, bool) {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil && n >= 0 {
		if n > int64((1<<63-1)/int64(time.Second)) {
			return time.Duration(1<<63 - 1), true
		}
		return time.Duration(n) * time.Second, true
	}
	if len(s) > 0 && strings.Trim(s, "0123456789") == "" {
		return time.Duration(1<<63 - 1), true
	}
	if t, err := http.ParseTime(s); err == nil {
		d := t.Sub(now)
		if d < 0 {
			d = 0
		}
		return d, true
	}
	return 0, false
}
func wait(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

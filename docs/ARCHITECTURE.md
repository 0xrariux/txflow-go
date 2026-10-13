# Architecture

Status: design for a community prototype, reviewed 2026-10-13. See [surface audit](TXFLOW_DEVELOPER_SURFACE.md) for why no live adapter exists.

## Packages and dependencies

```
root client/config -> markets, account, stream, trading, transport
markets/account/trading/orderbook -> types -> decimal
transport -> sdkerr
stream -> types, sdkerr
signing -> context (interfaces only)
demo -> root + embedded synthetic fixture
```

Public `Client` groups `Markets`, `Account`, `Stream`, and `Trading`. Market/account services delegate to small read interfaces. An absent adapter returns `ErrPendingConfirmation` before any I/O. `demo.NewClient` explicitly injects a synthetic backend. No implicit fallback to demo data is allowed. `WithNetwork` records a requested environment but resolves no URLs. `WithBaseURL` and `WithWebSocketURL` validate custom future configuration; they never fabricate endpoint routing. Generic `Transport` is independently testable; it is not wired to a TxFlow service yet.

Module path: `github.com/0xrariux/txflow-go`. This is a community repository with no TxFlow affiliation or endorsement.

| Package | Responsibility | Current boundary |
|---|---|---|
| `decimal` | Immutable finite decimal values, exact comparison/add/sub/halving, string JSON | Not a risk engine or arbitrary-division API |
| `types` | Proposed market, asset, book, trade, account domain models | JSON is SDK/demo schema, not verified TxFlow wire JSON |
| `markets`, `account` | Context-aware typed read services and validation | Adapter absent by default |
| `transport` | Bounded JSON HTTP, error types, retries, interval limiter, telemetry | No exchange endpoints or error-body mapping assumed |
| `orderbook` | Atomic snapshots and opt-in synthetic sequenced deltas | TxFlow sequence model unknown |
| `stream` | Future connector/session interfaces; bounded local replay | No WebSocket network client or reconnect implementation |
| `trading`, `signing` | Proposed requests and signer boundary | Fail closed / interfaces only |
| `demo` | Embedded authored fixtures and deterministic replay | Zero external connectivity |

## Financial precision decision (ADR-001)

Binary floats are unsuitable for price/quantity equality. `shopspring/decimal` is a candidate for production with broad decimal operations; `math/big.Rat` is exact but has no finite-decimal domain boundary; integer fixed scales would presume unverified asset scales. The prototype chooses a deliberately narrow standard-library-only decimal wrapper: canonical immutable strings, math/big arithmetic in temporary values, finite input notation, bounded input length, and exact halving. It has no generic division, rounding or exponent support. This keeps all tests/builds offline and avoids an early dependency commitment. Tests cover large values, fractional midpoint, malformed input and race safety. Re-evaluate a mature decimal library once supported scales and rounding rules are confirmed. These parser limits are SDK limits, never TxFlow limits.

## HTTP behavior

One context deadline bounds the whole operation, including throttling and retries; custom HTTP clients are copied without mutation. HTTPS is required except loopback HTTP for local tests. Relative request paths cannot escape the origin. Redirects are rejected. Default user-agent identifies `txflow-go`; no URL query, body, address, credentials or error body is logged by the telemetry callback.

Retries are off by default; at most five configurable retries. GET reads and explicitly marked read-only POST requests can retry transient network errors and 429/502/503/504 responses. Unsafe requests never retry. Backoff doubles to a cap with jitter; Retry-After is honored without retrying earlier than the server requests, and is bounded by the operation deadline. Oversized retry advice fails instead of being shortened. Application-level TxFlow error envelopes await a schema; current APIError means HTTP error only. Request-ID headers are best-effort generic HTTP metadata.

A context-aware shared minimum-interval limiter optionally limits all attempts. It makes no assumptions about TxFlow request weights. Hooks report method, attempt, status, elapsed time and coarse outcome, synchronously; hooks must be fast/concurrency-safe. No global logger and no Prometheus dependency.

## Order book and streams

Book state is mutex-protected, changes are validated before swapping, reads return copies. Duplicate prices, negative sizes and crossed books are rejected; zero-size deletion applies only to synthetic deltas. Depth and total levels are bounded. Strict contiguous sequence checking is explicitly selected for the local demo, never inferred from TxFlow. Gaps invalidate the book until a fresh snapshot; old/duplicate deltas return a typed stale error. A valid snapshot recovers state.

Local replay uses bounded channels and backpressure by waiting for a consumer or context cancellation. Cancellation/Close end its goroutine, close output channels and expose completion. There is no dropped-message strategy silently losing deltas. A future live connector must negotiate heartbeat, reconnect/resubscribe, snapshot buffering, duplicate identity and sequence rules before implementation. No ordering or user-stream guarantee is advertised.

## Proposed signing/trading boundary

Signer accepts an opaque, explicitly specified payload with format/version and returns a signature; no EIP-712 domain, hashing or keys are implemented. Trading methods return ErrPendingConfirmation. Nonce allocation, canonical encoding, wallet permissions and idempotency await official vectors. No endpoint submission exists.

## Test strategy

Offline tests use httptest and synthetic SDK JSON: cancellation, deadlines, rate limiting, retry bounds, Retry-After, unsafe POST, status/errors, malformed/trailing/null JSON, response limits, configuration, decimal precision, fixture parsing, book recovery and concurrent reads, replay cancellation. `go test ./...`, `go vet ./...`, `go test -race ./...`, formatting and CLI/example smoke checks are release gates. No live integration test exists; adding one requires written permission and a confirmed spec. Test coverage must describe actual behavior rather than mimic stubs.

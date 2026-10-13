# Technical summary for engineering review

Community prototype, 2026-10-13. This document describes current executable
behavior and the specific funded work needed to establish TxFlow compatibility.

## Evidence and scope

The official index leads to a Platform API placeholder and conceptual TIP SDK
interfaces. No public executable endpoint/signing/stream contract was established.
Published terms restrict automated access. The investigation did not proceed to
API probes or application bundle reverse engineering. A purported community SDK
was found as a secondary-source lead only; it is not a contract. See the full
surface matrix and source register for individual capability classifications.

Consequently, the prototype proves local engineering functionality rather than
live exchange behavior. This restraint is deliberate. Every dataset is authored
synthetic data; wire schema mapping, chain constants and signing vectors are absent.
Default clients return ErrPendingConfirmation. The demo adapter is explicitly
selected, so callers never silently receive fixtures when requesting live data.

## Public API and domain design

Client groups Markets, Account, Stream and Trading. Markets offers List, Assets,
Metadata, OrderBook, Trades and Ticker. Account offers State, OpenOrders,
HistoricalOrders, Fills and FundingHistory. Services accept contexts and return
proposed domain types. Small reader interfaces provide the adapter seam and must
be concurrency-safe. The default has no reader. WithNetwork does not invent
mainnet/testnet endpoints; custom base/WS options prepare future configuration.

Financial values are immutable finite decimal strings with temporary math/big
arithmetic for exact comparison, addition, subtraction and halving. There is no
float64 path, general division, risk calculation or rounding policy. The input
resource bound and accepted notation are SDK choices. Optional ticker fields
represent absent mark/oracle/funding/OI with nil rather than conflating absence
with zero. Current fixture JSON is SDK domain serialization, not TxFlow JSON.
A production decimal dependency remains a protocol-informed decision.

## HTTP infrastructure

The generic transport encodes request JSON and limits response bytes. One
operation deadline includes rate-limit waits, network attempts and backoff.
Custom clients are copied without mutation; redirects are rejected. HTTPS is
required except loopback HTTP for local tests. Relative paths stay on the chosen
origin. User-agent identifies the prototype.

Retries default to zero. Configurable attempts are bounded and safe GETs or
explicitly marked read-only POSTs may retry transient transport failures and
selected HTTP statuses. Trading submissions are never wired into this transport.
Exponential backoff has a cap/jitter; Retry-After is respected without retrying
early. Error bodies are not echoed to prevent accidental disclosure. APIError
currently maps HTTP statuses only; TxFlow business codes and envelopes await
engineering-approved schemas. Optional hooks expose sanitized method/attempt/
status/latency/outcome and the interval limiter makes no quota-weight claims.

## Book and stream mechanics

Local Book atomically swaps validated snapshots and delta results under a mutex.
Reads copy slices and sequence pointers. Price keys use canonical decimals;
duplicate levels, negative sizes and crossed books are invalid. Top-N depth and
exact best bid/ask, mid and spread are available. Memory is bounded by a level
cap. The demo selects strict contiguous synthetic sequences: duplicate events
are rejected, gaps invalidate state, and an authoritative new snapshot recovers.
No assertion is made that TxFlow provides this sequence model.

Streaming currently consists of Connector/Session interfaces and a finite local
ReplayConnector. Subscription channels are bounded; a slow consumer blocks the
producer until it reads or cancels. Close/context cancellation terminate the
producer and close channels. This is not a network WebSocket implementation.
Heartbeat, reconnect/resubscription, connection-state events, actual malformed
wire handling and account-stream authorization remain design requirements.

The key demonstration loads a synthetic market snapshot, consumes replay updates,
prints quotes/depth, rejects a duplicate, detects a gap and recovers at sequence
103 before continuing to 104. The CLI requires --demo for data; a normal client
returns a clear pending error. Examples are small and all operate offline.

## Security and signing

Signer is an interface with public address, explicitly identified payload format/
version/data and signature format/bytes. No ECDSA, private-key parser, EIP-712
domain, action hashing, serialization or nonce allocator is shipped. Trading's
proposed PlaceOrder/CancelOrder/CancelAll/ModifyOrder methods fail locally.
Neither calling code nor CLI can place a live order with this implementation.

Agent design separates proposed trading authority from owner transfer/withdrawal
authority and covers storage, rotation, expiry and revocation. These controls
must be enforced by the official protocol before they are advertised. Go cannot
promise complete zeroization of runtime key copies. Hardware/external/callback
signers are future implementations sharing canonical vectors, not separate
ad hoc encoding paths.

## Validation and production gates

Tests cover local HTTP headers/JSON, safe retry bounds and unsafe POST behavior,
Retry-After, cancellation/deadlines, origin/redirect constraints, malformed and
oversized responses, precise arithmetic, domain fixtures, account/market
validation, independent returned values, concurrent readers, book gaps/recovery,
replay cancellation and CLI source labels. No production services are used.
Race/vet/format/build checks and CI are repository gates; GitHub-hosted CI will
only be executable after publication. Local validation is recorded separately.

Before v1, engineering should supply supported environments, schemas, pagination/
retention, quotas/errors, sequence/lifecycle rules, canonical positive/negative
signing vectors and a permissioned testnet path. Stress/leak/soak tests, benchmarks,
verified user streams, action-status/idempotency and delegation scope tests are
funded continuation. A 0.1 prototype is not a production safety certification.

## Funding and collaboration

Preferred six-milestone scope: 220h at an indicative $125/hour, **$27,500**.
Alignment is first and can re-scope later deliverables if unsupported. Signer and
trading acceptance require official vectors and testnet lifecycle checks; stream
acceptance requires controlled failures and verified recovery. The public module is github.com/0xrariux/txflow-go. Future ownership and ongoing
maintenance should be agreed with collaborators; security reporting is described
in SECURITY.md. No grant award is assumed.

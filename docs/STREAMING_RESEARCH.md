# Streaming boundary and missing protocol details

Status: pending live TxFlow specification and automated-access permission (T2/T6 in [SOURCES](SOURCES.md)). No WebSocket endpoint, subscription envelope or heartbeat was verified. Therefore no network WS client or reconnect behavior is shipped.

Implemented: Connector/Session interfaces, explicit local ReplayConnector, context cancellation, graceful shutdown, bounded buffering and producer backpressure. Offline event schema is proposed SDK domain data. Synthetic replay demonstrates duplicate rejection, contiguous sequence gaps and snapshot recovery in the separate orderbook package. These are **local test policies**, never asserted TxFlow guarantees.

Required before live implementation:

| Concern | Required evidence | Planned behavior |
|---|---|---|
| URL, environments and lifecycle | Official supported hosts/TLS/connection caps | Explicit approved environment configuration |
| Subscription acknowledgments | Request/ack/error schemas and subscription IDs | Register only confirmed subscriptions; typed errors |
| Heartbeat | Ping/pong framing and idle timeout | Context-aware liveness supervision |
| Reconnect/resubscribe | Retry guidance, auth renewal, snapshot validity | Capped jittered backoff and restoration events |
| Books | Snapshot vs delta, sequence scope, buffering boundary | Resume only with proven continuity; invalidate on gaps |
| Duplicate events | Stable event identity and replay windows | Bounded dedup if supported |
| Backpressure | Full snapshot or delta delivery semantics | Fail/resync rather than silently drop deltas |
| Account streams | Auth/user scope and authorization | Isolated subscriptions and telemetry without secrets |
| Malformed messages | Wire schema and semantic error rules | Size bounds, typed decode errors and recovery policy |
| Observability | Supported state transitions | Optional reconnect/received/dropped/limit hooks |

Planned events: Connecting, Connected, Reconnecting, Resynchronizing, Closed. Current replay does not emit connection events or reconnection metrics because it has no network lifecycle. Testnet recovery/heartbeat/reconnect tests are grant scope. A custom WebSocket URL option is configuration preparation; it does not activate networking.

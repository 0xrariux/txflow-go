# TxFlow developer surface audit

Reviewed 2026-10-13. Evidence IDs resolve in [SOURCES](SOURCES.md).

## Decision

The API reference is a placeholder (T2). Product and TIP docs establish concepts but no executable integration contract (T1/T3/T4). T6 expressly restricts automated access. This prototype therefore ships **no TxFlow endpoint mapping, production network default, API probes, signer or trading submission**. Offline synthetic fixtures and generic transport demonstrate engineering feasibility. Live functionality needs both an authorized contract and permission.

Classification: `OFFICIAL_DOCUMENTED` means the specific claim is stated in official documentation; `OFFICIAL_OBSERVED` means observed from an official interface without a stability promise; `COMMUNITY_OBSERVED` means third-party behavior only; `INFERRED` means a hypothesis; `UNKNOWN` means insufficient evidence. Confidence refers to the claim, not presumed endpoint correctness.

## Product evidence, separate from API evidence

| Capability | Classification | Evidence | Confidence | Suitability | Confirmation |
|---|---|---|---|---|---|
| Perpetual product exists | OFFICIAL_DOCUMENTED | T7 | High | Domain vocabulary only | Yes, query schema |
| Margin/account concepts | OFFICIAL_DOCUMENTED | T8/T4 | High | Proposed domain vocabulary | Yes, query schema |
| Prediction documentation exists | OFFICIAL_DOCUMENTED | T9 | High | Research; provider boundary unresolved | Yes |
| TIP SDK concept | OFFICIAL_DOCUMENTED | T3 | High | Architecture inspiration | Yes, actual distribution |
| App HTML accessible | OFFICIAL_OBSERVED | T10 | High | No API permission inferred | Yes |
| Community trading SDK lead | COMMUNITY_OBSERVED | C1 secondary | Low | Investigate only | Yes |
| Hyperliquid-like protocol hypothesis | INFERRED | C1 lead | Low | Do not implement from resemblance | Yes |

## API implementation matrix

Every row concerns the **wire contract**, even when a corresponding product concept is documented. Unknown does not mean unsupported. No speculative endpoint has been promoted to official status.

### Infrastructure

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| base API endpoints | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| environments | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| testnet/mainnet | UNKNOWN | T1 indexes onboarding guides; API addresses unconfirmed | Low for API | No live implementation | Yes |
| versioning | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| authentication | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| rate limits | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| errors | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| request IDs | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| pagination | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### Markets

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| market list | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| asset metadata | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| spot markets | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| perpetual markets | UNKNOWN | T4/T7 describe concepts; T2 supplies no query contract | Low for API | No live implementation | Yes |
| prediction markets | UNKNOWN | T9 references Probly; API/provider boundary unconfirmed | Low for API | No live implementation | Yes |
| prices | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| tick sizes | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| lot sizes | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| funding | UNKNOWN | T4/T7 describe concepts; T2 supplies no query contract | Low for API | No live implementation | Yes |
| mark price | UNKNOWN | T4/T7 describe concepts; T2 supplies no query contract | Low for API | No live implementation | Yes |
| oracle price | UNKNOWN | T4/T7 describe concepts; T2 supplies no query contract | Low for API | No live implementation | Yes |
| open interest | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### Order book

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| snapshot | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| incremental updates | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| depth | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| sequence numbers | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| resynchronization | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### Trades

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| recent trades | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| historical trades | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| user fills | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### Accounts

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| balances | UNKNOWN | T4/T8 describe concepts; T2 supplies no wire schema | Low for API | No live implementation | Yes |
| positions | UNKNOWN | T4/T8 describe concepts; T2 supplies no wire schema | Low for API | No live implementation | Yes |
| margin | UNKNOWN | T4/T8 describe concepts; T2 supplies no wire schema | Low for API | No live implementation | Yes |
| PnL | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| open orders | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| trade history | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| funding history | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### WebSocket

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| URL | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| connection lifecycle | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| subscriptions | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| heartbeat | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| reconnect | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| sequence guarantees | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| market streams | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| account streams | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### Trading

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| place order | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| cancel | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| cancel all | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| modify | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| batch | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| leverage | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| margin mode | UNKNOWN | T4/T8 describe concepts; T2 supplies no wire schema | Low for API | No live implementation | Yes |
| transfers | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| withdrawals | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

### Authentication / signing

| Capability | Classification | Evidence / notes | Confidence | Suitable now | TxFlow confirmation |
|---|---|---|---|---|---|
| wallet authentication | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| agents | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| EIP-712 | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| nonce semantics | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| serialization | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| action hashing | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |
| replay protection | UNKNOWN | T2: no public contract; T6: access alignment required | Low for API | No live implementation | Yes |

## Practical gaps and blockers

Critical: permission for automated integrations and authoritative schemas for read queries, streams and signing. High: fixtures/testnet access, sequence recovery contracts, documented errors/limits and version policy. Medium: official SDK discovery, examples and compatibility CI. A native Go layer is a proposed opportunity, not an established claim that TxFlow has no SDK.

No financial balances, tick sizes, lot sizes or stream envelopes were captured. Every shipped dataset is authored synthetic data. The SDK JSON used in demo files is its own domain serialization, never claimed to be TxFlow JSON. A live order book demonstration cannot be delivered responsibly from this evidence; the local replay demonstrates mechanics separately.

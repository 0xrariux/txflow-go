# Trading SDK benchmark

As of 2026-10-13. Evidence IDs link to [SOURCES](SOURCES.md). This is a bounded documentation review, not a runtime benchmark or ranking of exchanges. “Not established” is a research limitation. A source directory is not proof that all its tests pass. Stability/semantic versioning are not assumed from package version numbers.

## SDK availability and ownership

| Ecosystem | Go | Python | TypeScript | Evidence |
|---|---|---|---|---|
| TxFlow | No official package established | No official package established | TIP SDK concept; community lead unverified | T2/T3/C1 |
| Hyperliquid | Official Go not established by reviewed API index | Linked official Python SDK | Linked community SDKs | H1 |
| dYdX | Go trading client not established in reviewed client monorepo | Python v2, maintained by Nethermind | Chain client and lower-level dydxjs | D1/D2 |
| Drift | Native Go client not established | DriftPy | Protocol repository SDK | R1/R2; both redirected/archived at review time |
| Coinbase Advanced | Coinbase sample Go REST client | Official SDK | Sample SDK | B1/B2 |
| Paradex | Official Go not established | Primary paradex-py client | Official TypeScript not established by reviewed sources | P1/P2 |

## Developer infrastructure

| Ecosystem | REST / streaming quality signals | Signing abstraction | Examples / bots / market making | Errors / rate limits | Stability / semver / CI evidence |
|---|---|---|---|---|---|
| TxFlow | Platform API placeholder; TIP interfaces conceptual | No canonical executable spec established | Public community lead suggests interest; compatibility unverified | No official wire model or numeric quota established | No public API compatibility guarantee established |
| Hyperliquid | API entry point and environment conventions documented | SDK is a reference implementation; TxFlow cannot inherit it | SDK links lower initial integration work | Weighted REST and WS limits explicit | Must pin SDK versions and inspect compatibility; no guarantee inferred |
| dYdX | Integration docs and distinct chain clients | dydxjs handles chain message encoding | Multiple client languages; grant-funded third-party C++ precedent | Full error/quota review deferred | Client monorepo exposes CI directory; semver contract not audited |
| Drift | RPC/account subscriptions rather than assuming exchange REST equivalence | Chain client manages wallet transactions | Python examples and provider-specific subscription setup | RPC-provider limits affect integrations | Historical repos include tests/workflows; archive status requires maintenance diligence |
| Coinbase Advanced | REST and WS clearly separated | SDKs wrap authenticated requests | Go sample gives a concrete backend starting point | Full endpoint-specific quota/error audit deferred | Official vs sample lifecycle distinction explicit; release policy not audited |
| Paradex | Docs expose REST/WS and Python client | Account/subkey abstractions in client | Async/threaded WS and runnable examples | Full quota audit deferred | Primary repo includes tests/workflows; no stability guarantee inferred |

These signals are actionable design patterns: explicit ownership, concrete install paths, transport-specific documentation, signer isolation and visible tests. They do not imply permission to automate TxFlow or compatibility between protocols.

## Developer Tooling Gaps TxFlow Could Address

| Rank | Opportunity | Acceptance evidence |
|---|---|---|
| Critical | Explicit third-party automation authorization and supported API contract | Written permission and versioned schema reviewed by engineering |
| Critical | Canonical signing/nonces/replay specification | Official positive and negative cross-language vectors |
| High | Snapshot/stream continuity, backpressure and recovery guidance | Defined sequence semantics, heartbeat and testnet recovery tests |
| High | Machine-readable errors and quotas | Typed codes, limit scopes, weights and Retry-After behavior |
| High | Native Go implementation aligned with supported backend workflows | Approved adapter, public account/market fixtures and integration CI |
| Medium | SDK discoverability, ownership and release policy | Official links distinguish supported and community tooling |
| Medium | Bot and market-maker starter projects | Testnet-only examples with explicit authorization boundaries |
| Nice-to-have | Metrics exporters and advanced tooling | Demand demonstrated after core SDK adoption |

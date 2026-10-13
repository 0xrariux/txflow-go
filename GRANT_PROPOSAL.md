# TxFlow Go SDK — Ecosystem Infrastructure Proposal

Draft for engineering/ecosystem review, 2026-10-13. Community project; no TxFlow endorsement, grant award or official compatibility is claimed.

## Executive Summary

This proposal offers a native Go developer layer for services integrating with TxFlow. Reviewed official documentation describes TIP SDK/interface concepts, but the [Platform API reference](https://docs.txflow.com/api-reference/readme) currently says Coming Soon. The opportunity is to collaborate on a supported Go integration, rather than assume undocumented application behavior is stable. Existing or planned official SDK efforts should be checked with the team before funding a parallel effort.

A self-funded prototype already demonstrates precise financial models, context-aware typed market/account services, reusable HTTP infrastructure, a CLI, offline fixtures, bounded stream replay and local book recovery. These are functional **offline engineering foundations**, not a live exchange adapter. The [evidence register](docs/SOURCES.md) records the investigation; the [terms](https://docs.txflow.com/terms/terms-of-use) make express written automated-access permission a first milestone dependency.

Funding would turn those foundations into an authorized, tested integration: canonical query mappings and signing, permissioned testnet trading, resilient WebSockets, confirmed agent-wallet support, hardening and a public release. The preferred discussion scope is **$27,500**, with milestone acceptance and final scope aligned with TxFlow. If supported APIs or delegation differ from this proposal, the deliverables should be revised before implementation.

## Problem

Backend developers otherwise need to choose wire formats, exact numeric handling, transport policies and stream recovery individually. An SDK can share that work across bots, market makers, indexers and monitoring services. The research surfaced a community SDK lead, but did not validate its source or compatibility; it is cautious evidence of interest, not proof of market demand. No unsupported estimate of users, volume or integration savings is made.

## Proposed Solution

An independently authored native Go SDK with a small typed API, separate verified wire adapters, signer isolation, explicit compatibility policy and realistic bot examples. Public/read-only functionality comes first. Authenticated operations ship only after permission and protocol validation. The production target includes durable market/account reads, precise orders, recovery-aware streams and clear errors/quotas.

## Why Go?

Go suits long-running services through contexts, goroutines, networking libraries, race detection and simple deployment. The SDK should let operators compose concurrent market/account work and integrate existing telemetry without imposing a framework. It complements other languages rather than displacing them.

## Current Prototype

| Delivered | Evidence / limitation |
|---|---|
| Surface audit and six-ecosystem benchmark | docs/SOURCES.md and SDK_BENCHMARK.md; bounded research |
| Architecture and typed services | Explicit domain-vs-wire boundary; no endpoint mapping |
| Exact decimal values | Standard-library narrow operations; tested precision |
| HTTP transport and configuration | Deadlines, response bounds, safe opt-in retries, Retry-After, limiter, hooks |
| Market metadata/books/trades/ticker and account query UX | Authored synthetic fixtures only |
| Local stream and order-book core | Cancellation, backpressure, duplicates, synthetic sequence gaps/recovery |
| CLI and five short examples | Explicit offline source labels; no live trading commands |
| Signing/trading/agent design | Interfaces and local pending errors only |
| Offline tests and repository documentation | Meaningful transport, fixtures, validation and concurrency checks |

The prototype covers a modest foundation toward the envisioned SDK. A precise 20–35% completion percentage cannot be justified before the real protocol scope is agreed; most exchange-specific work remains unfunded and unimplemented.

## Grant Scope

Align permitted access and contracts; implement official read mappings; verify canonical signing; add supported trading actions and error/status handling; implement real WS recovery and account streams; add confirmed delegation; harden and release. No future FIX gateway, strategy framework, indexer or MCP server is included in the preferred scope. No mainnet-value trading is required for acceptance.

## Deliverables and acceptance

| Milestone | Deliverable | Reviewable acceptance | Effort / budget |
|---|---|---|---|
| 0 — Self-funded prototype | Existing offline foundations and evidence package | Local reproducible demo and passing offline checks; live gaps clearly labeled | Already delivered; $0 requested |
| 1 — Official API alignment | Written access permission, supported hosts, read schemas, signed-action vectors, quotas/errors, API compatibility matrix | Engineering review; authorized market/account fixtures; first verified read adapter and testnet query smoke checks | 20h / $2,500 |
| 2 — Authenticated trading | Canonical payloads/signer, supported order/batch/cancel/modify actions and leverage/margin controls | Official positive/negative cross-language vectors; permissioned testnet lifecycle tests; ambiguous-outcome handling; no unsafe automatic retries | 60h / $7,500 |
| 3 — Production WebSocket | Supported market/user streams, heartbeat/reconnect/resubscribe, bounded buffering and local book | Fault-injection recovery tests and verified sequence policy; no silent delta loss; testnet soak | 50h / $6,250 |
| 4 — Agent infrastructure | Delegation lifecycle and isolated signer, only if officially supported | Permission tests for scope/expiry/revocation/rotation and withdrawal separation | 20h / $2,500 |
| 5 — Hardening | Confirmed limiter/errors, stress/race/leak tests, benchmarks and stream/request telemetry | Published test results, finite resource bounds, documented failure/recovery limits | 40h / $5,000 |
| 6 — Ecosystem release | Owner/module resolution, version policy, docs, bot/market-maker starter examples and release workflow | Reviewed v1.0 candidate; reproducible examples; compatibility matrix and maintenance handoff | 30h / $3,750 |

Total: **220 hours / $27,500** at an indicative $125/hour planning rate. These are effort estimates, not external price claims. A single contributor at approximately 25 hours/week implies about nine working weeks after prerequisite access and timely engineering feedback. Protocol delays are outside that estimate. M2 and M3 may proceed after M1; M4 requires a supported authorization model; M5/M6 depend on accepted functionality.

If M1 reveals incompatible scope or absent features, stop affected live work and agree a narrower scope or reassigned milestone. Payments should follow accepted deliverables; no unsupported feature should be represented as delivered. A proposed 30-day release stabilization window is included in M6's allocated effort; longer ongoing maintenance needs a separately agreed retainer or community ownership plan.

## Governance and ownership

MIT is selected for broad reuse of original implementation; engineering can discuss contribution ownership without implying endorsement. The community repository is [0xrariux/txflow-go](https://github.com/0xrariux/txflow-go); private security reporting is described in [SECURITY.md](SECURITY.md). Future ownership and maintenance can be aligned with TxFlow. No third-party signing code or proprietary application bundles are included. [Budget alternatives](docs/GRANT_BUDGET.md), [roadmap](ROADMAP.md), [questions](docs/QUESTIONS_FOR_TXFLOW.md) and [contact package](proposal/ONE_PAGER.md) support this discussion.

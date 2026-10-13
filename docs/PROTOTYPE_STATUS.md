# Prototype status

Version: **v0.1.0-prototype**. Community project, not affiliated with or endorsed by TxFlow.

## Implemented

- Exact finite decimal values and typed market/account read services.
- Generic HTTP JSON transport with context deadlines, bounded responses, typed
  HTTP errors, opt-in safe-read retries, Retry-After, limiter and telemetry.
- Explicit synthetic fixtures, local stream replay and order-book recovery.
- Read-only CLI and five short examples.
- Signing/trading interfaces and protocol research; trading methods fail locally.

All demonstration data is synthetic. There is no live TxFlow adapter or implicit
fallback from a live request to fixtures. Domain JSON is not a verified wire schema.

## Validation

The Go checks workflow runs unit tests, vet, builds, formatting, documentation
links and the offline demo on Go 1.23 and stable. Stable also runs race tests.
Current hosted results are available in [GitHub Actions](https://github.com/0xrariux/txflow-go/actions/workflows/ci.yml).

Local tests, race tests, vet, builds, examples and precision fuzzing passed with
Go 1.26.5 on macOS/arm64. Tests use local HTTP servers and authored fixtures;
no external service or private key is needed.

## Pending TxFlow alignment

Live read APIs, WebSocket lifecycle and sequence recovery, signing/nonces/replay,
agent permissions, authenticated trading, quotas/errors and API compatibility.
Written automated-access permission and canonical contracts/vectors are prerequisites.
See the [surface audit](TXFLOW_DEVELOPER_SURFACE.md), [signing research](SIGNING_RESEARCH.md)
and [engineering questions](QUESTIONS_FOR_TXFLOW.md).

## Proposed continuation

The [grant proposal](../GRANT_PROPOSAL.md) describes conditional milestones.
The preferred discussion scope is $27,500; no award or official SDK status is claimed.

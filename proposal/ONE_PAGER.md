# txflow-go — Native Go infrastructure for TxFlow

**Community proposal · 2026-10-13 · no official endorsement**

**Problem.** Backend builders benefit from a shared integration layer for
precise market data, account reads, signing and recovery. Reviewed official
TxFlow docs describe TIP components; the Platform API page is Coming Soon.
Existing/planned official SDK work should be checked before funding overlap.

**Solution.** A native Go SDK aligned with supported contracts, small typed
services and separate transport/signing boundaries. Written automated-access
permission is an explicit prerequisite under the published terms.

**Existing prototype.** Exact decimals, generic HTTP deadlines/retries/errors,
optional throttling/hooks, typed market/account query interfaces, synthetic
fixtures, bounded local replay, gap-aware local book, CLI, five short examples,
offline tests and protocol/grant documentation. **No live adapter, real WS,
signer or trading submitter is implemented.**

**Value.** Shared backend infrastructure for market makers, bots, data services
and agents; a common verified vector suite can reduce repeated protocol work.
No adoption, volume or performance figures are invented.

**Requested support.** Discuss **$27,500** across six accepted milestones:
API alignment/reads, authenticated trading, real streaming, confirmed agents,
hardening and release. Core $17,500 and Extended $45,000 alternatives are
available; final amount and ownership should be agreed with TxFlow.

**Next step.** A short engineering/ecosystem review of access permission,
schemas/vectors, existing SDK plans and milestone fit. The demo runs locally:

```sh
go run ./cmd/txflowctl --demo stream book BTC-USDC
```

[Technical summary](TECHNICAL_SUMMARY.md) · [Grant request](GRANT_REQUEST.md) ·
[Full proposal](../GRANT_PROPOSAL.md) · [Engineering questions](../docs/QUESTIONS_FOR_TXFLOW.md).

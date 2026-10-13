# Grant budget alternatives

Planning estimates as of 2026-10-13, in USD. Final scope, amount and payment denomination must be aligned with TxFlow. No grant program or award is assumed. Estimates use $125/hour to make effort transparent; this is a proposal assumption, not a market-rate finding. Prototype work is self-funded and excluded from requested amounts.

| Milestone | Core hours / USD | Full hours / USD |
|---|---|---|
| 1: permission, schemas, verified reads | 20 / $2,500 | 20 / $2,500 |
| 2: supported authenticated trading | 40 / $5,000 | 60 / $7,500 |
| 3: real streaming | 25 / $3,125 | 50 / $6,250 |
| 4: confirmed agent infrastructure | 0 / $0 | 20 / $2,500 |
| 5: hardening | 35 / $4,375 | 40 / $5,000 |
| 6: release/docs/stabilization | 20 / $2,500 | 30 / $3,750 |
| Total | **140 / $17,500** | **220 / $27,500** |

## Core — indicative $15k–$20k

Minimum supported production SDK, scoped to confirmed market/account reads,
basic limit-order place/cancel, public market streams, exact numbers and
essential errors/limits/tests/docs. Excludes advanced modify/batch breadth,
user-stream expansion, delegation lifecycle and market-maker starter depth.
If this breadth exceeds 140h after alignment, reduce it before acceptance.
Production readiness requires meeting the agreed fault/recovery gates;
smaller budget must not imply skipped signing validation.

## Full — indicative $25k–$35k

Preferred **$27,500** (discussion band $25k–$30k). Covers all six funded
milestones in [GRANT_PROPOSAL](../GRANT_PROPOSAL.md): supported trading breadth,
user streams, agent lifecycle if supported, hardening and starter examples.
Budget is contingent on obtaining official contracts and permission. Additional
protocol ambiguity should prompt scope negotiation, not unsupported inference.

## Extended — indicative $40k–$50k

Illustrative **$45,000 / 360h**: Full 220h plus 40h advanced bot tooling
($5,000), 50h market-maker starter infrastructure ($6,250), 30h additional
benchmarks/fault simulation ($3,750), and 20h extended maintenance ($2,500).
No trading strategy profitability, liquidity commitment, FIX gateway or
indexer service is promised. These extras require demonstrated ecosystem demand.

## Delivery and exclusions

Full estimate assumes roughly nine working weeks at 25h/week, after access
and timely engineering reviews. Core roughly six; Extended roughly fifteen.
These are effort-based planning horizons, not fixed-date commitments. Milestone
acceptance governs payment. Unsupported agent functions must be removed or
re-scoped explicitly. Mainnet capital, exchange fees, third-party hosting/KMS
costs and separate external security audit are not included. No hidden
contingency or invented adoption/volume figures are added. A 30-day release
stabilization window consumes part of M6, not unlimited maintenance hours.

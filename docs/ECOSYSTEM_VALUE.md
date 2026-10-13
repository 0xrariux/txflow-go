# Ecosystem value

Proposed value, not measured adoption. This assumes TxFlow authorizes the
integration and supplies supported contracts. [SOURCES](SOURCES.md) records
what has actually been established.

| Builder | Concrete benefit | Prototype proof | Funded continuation |
|---|---|---|---|
| Market makers | Shared precise book/quote/recovery primitives | Synthetic depth, mid/spread and gap recovery | Live book continuity and testnet order lifecycle |
| Trading bots | Typed context-aware services and isolated signing | Query interfaces, transport and explicit pending boundary | Canonical signer, supported orders and rate limits |
| Exchanges/custodians | Reviewable backend integration contracts | Exact domain models and HTTP errors | Official compatibility matrix, authority isolation |
| Data infrastructure | Consistent market/account parsing | Domain fixtures and metadata/trade services | Retention/pagination contract and live indexing inputs |
| AI agents | Programmatic read interface without exposing keys | Keyless synthetic reads and offline CLI | Confirmed delegated trading scopes if supported |
| Builders | Reproducible first integration example | Five short examples, no network setup needed | Approved install path and permissioned testnet examples |
| TxFlow engineering | One shared implementation and vector suite | Evidence discipline and concise alignment questions | Reduced duplicated protocol support; reviewed releases |

Measure impact after release using opt-in reports: time to first successful
query, time to first testnet order, number of supported integrations, reproducible
bug reports, recovery success during injected failures and repeat usage.
Do not invent baseline savings or collect wallet secrets to measure adoption.
The current demo proves local engineering functionality only; it does not prove
exchange latency, bot profitability, liquidity improvement or production compatibility.

# Questions for TxFlow engineering and ecosystem teams

Prioritized from the 2026-10-13 [surface audit](TXFLOW_DEVELOPER_SURFACE.md).
These questions are prerequisites, not claims about supported features.

1. Can TxFlow give express written permission for a third-party Go SDK's read-only probes, testnet automation and later trading integration, consistent with Terms 3.1(l)/(o)? What scopes and restrictions should apply?
2. Is a public or partner Platform API spec available while the official page says Coming Soon? Which hosts, environments, endpoints and market modules are intended stable third-party contracts?
3. What are the canonical query schemas, asset identifiers, decimal/rounding rules, pagination and historical retention limits? Can engineering provide sanitized market/account fixtures?
4. Which public and account WebSocket streams are supported, with what subscribe/ack/error/heartbeat lifecycle and connection/subscription limits?
5. Are books snapshots or deltas? What sequence scope, duplicate identity, snapshot buffering boundary and resynchronization procedure should clients use?
6. Can engineering provide the complete signing specification and positive/negative vectors: domains, action types, canonical bytes, hashes, nonces, expiry and replay protection?
7. Are agent wallets supported for third-party bots? What approval, scope, expiry, rotation and revocation rules distinguish trading from transfer/withdrawal authority?
8. Which trading actions and safeguards are supported: batch/cancel-all/modify, leverage/margin, client order IDs, status queries and idempotency after ambiguous submissions?
9. What are official quotas, retry guidance, typed errors, request IDs and API version/deprecation guarantees? Is there a supported testnet CI integration path?
10. Is a Go SDK already planned or available, and would TxFlow prefer community or organization ownership? Is a milestone-funded contribution mechanism available, and which of the proposed scopes best fits the team's priorities?

First engagement should resolve 1–3 before live work. No SDK ownership,
grant program, response time or acceptance is assumed.

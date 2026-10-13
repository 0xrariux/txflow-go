# Signing research

Research date: 2026-10-13. Evidence classifications refer to [SOURCES](SOURCES.md).
No signatures were generated, no signed probes submitted, and no signing
algorithm inferred from another venue was implemented.

## Confirmed by official TxFlow documentation

T2 is a placeholder Platform API page. T3/T4 describe SDK/interface concepts
without canonical signing bytes. T6 restricts automated access. This establishes
a permission/specification gate, not a signing contract. No official executable
EIP-712 domain, chain ID, action hash or nonce rule was established by the
reviewed sources.

## Observed from public TxFlow interfaces

The public app HTML exists (T10), but was not used to infer signing. Bundles
were not downloaded. No wallet challenge, nonce, signature, delegation,
withdrawal or transfer payload was captured. Official observed signing
evidence: none.

## Community implementations

C1 is a secondary mirror lead for a purported trading SDK. Its underlying
primary implementation, version availability and correctness were not
verified. It motivates an engineering question but supplies no adopted
algorithm, chain constant, vector or security guarantee. X1/X2 are unrelated
name collisions and were excluded. No claim of official SDK absence follows
from these searches.

## Unverified assumptions

| Question | Known today | Required evidence |
|---|---|---|
| EIP-712 domain | Unknown | Name, version, chainId, verifyingContract, salt and environment mapping |
| Action structures | Unknown | Order/cancel/modify/batch/user-action schemas and exact field types |
| Nonces | Unknown | Scope, units, windows, monotonicity, concurrent allocation and expiration |
| Serialization | Unknown | Field order, numeric/text encoding, zero normalization and optional fields |
| Hashing | Unknown | Canonical inputs, intermediate hashes and signature recovery vectors |
| Replay protection | Unknown | Chain/version/user/agent binding and duplicate submission semantics |
| Agent delegation | Unknown | Approval schema, allowed scopes, expiration, rotation and revocation |
| Withdrawals/transfers | Unknown | Separate owner authority, destination validation and signature rules |

Resemblance to Hyperliquid is not evidence of compatibility. In particular,
MessagePack, phantom-agent structures, Ethereum chain IDs or particular
signature layouts must not be imported as TxFlow facts.

## Proposed implementation architecture

A spec-versioned payload builder will own canonical encoding and validation.
A nonce provider will coordinate per verified scope. Signer receives explicit
format/version/data and returns an explicitly typed signature; it should not
know HTTP or retry behavior. A submitter will verify action authorization and
handle ambiguous network outcomes via official status/idempotency APIs.
Unsafe write retries will remain disabled until that contract is established.

Local ECDSA, external/hardware, callback and agent-wallet signers should share
the same vector suite. They remain options, not existing implementations.
The current narrow interface avoids pulling an Ethereum stack into an offline
prototype; it may change once schemas and address semantics are confirmed.

No test vector is fabricated. Milestone 1 must obtain positive and negative
vectors: domain mismatch, reordered fields, trailing zeros, optional vault,
large nonce, unauthorized action and replay. Milestone 2 acceptance requires
cross-language equivalence and permissioned testnet-only submission. Production
keys and live-value transactions are outside prototype validation.

# Agent wallet design

**Proposed security model, not verified TxFlow behavior.**
No official agent delegation contract was established (T2/T3 in
[SOURCES](SOURCES.md)). Community lead C1 does not establish permissions.

## Purpose and authority

A backend trading signer should have narrowly scoped delegated authority and
should not hold the wallet used for deposits or withdrawals. TxFlow engineering
must confirm whether that separation is supported and enforced on chain.
Until confirmed, this document is design guidance only. Market/account reads
should remain keyless where an approved public contract permits them.

## Proposed Go API

An eventual separate authorization service could expose ApproveAgent,
InspectAuthorization, RevokeAgent and RotateAgent, with explicit scope,
expiry and principal. Owner and agent Signer objects must be distinct.
No such TxFlow actions or methods are implemented in this prototype.
Approval, transfer and withdrawal signing should never be implicit during
a bot's startup or an order request.

## Key lifecycle

Use external signers/KMS/hardware where practical. An in-memory signer must
receive a key explicitly, avoid logs and keep access short-lived; Go cannot
promise complete zeroization of all runtime copies. Avoid plain-text config,
source code, shell history and committed environment files. Secrets should
not enter telemetry hooks, exception text or CLI output.

Provision a separate key per confirmed authorization/nonce isolation scope.
Rotate by approving the replacement, confirming scopes and expiration,
draining pending requests, revoking the old key and verifying revocation.
Exact ordering and overlap rules need TxFlow confirmation. Do not assume
revocation invalidates already signed actions or historical nonce state.
If the platform does not enforce least privilege, document that limitation
and avoid presenting an agent as withdrawal-safe.

## Failure model and acceptance

Bot compromise should expose only the permitted trading scope, with explicit
order/value/rate limits if supported. Owner compromise remains a separate risk.
No unsupported permission controls may be promised. Testnet acceptance must
include unauthorized transfer/withdrawal rejection, expiration, revocation,
rotation under concurrent orders and replay after delegation changes.
Canonical schemas and written automation permission are prerequisites.

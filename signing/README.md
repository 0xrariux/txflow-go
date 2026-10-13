# Signing boundary

Interfaces only. No private-key parser, ECDSA, EIP-712 domain, hashing,
nonce allocator, canonical serialization or action submitter is implemented.

Payload and signature formats are proposed Go abstractions. They are not
TxFlow wire types or signing vectors. Trading operations never call Signer.

Production implementation needs the [research gates](../docs/SIGNING_RESEARCH.md):
official action/domain schemas, canonical byte vectors, replay rules, network
separation, allowed agent scopes and written automated-access permission.
Local ECDSA, external/hardware, callback and agent signers are future choices,
subject to the same approved payload builder and vector suite.

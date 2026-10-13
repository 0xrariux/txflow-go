# Security

## Reporting a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/0xrariux/txflow-go/security/advisories/new)
to report a vulnerability confidentially. Include the affected version,
reproduction steps and expected impact. Do not include private keys, seed phrases,
credentials or sensitive wallet data. Avoid public issues for security reports.

## Supported version

The v0.1.0-prototype prerelease is maintained as an offline prototype. No response
SLA or production compatibility guarantee is offered.

## Security boundaries

This prototype has no live trading adapter, signer, private-key parser, delegation
client, deposit/transfer/withdrawal API or production compatibility guarantee.
It uses explicit synthetic fixtures for demonstrations. Signing and trading
research documents the permission/schema/vector gates for future work.

Transport hooks omit request URLs, bodies and addresses. HTTP errors do not echo
server bodies. Synthetic fixtures are not signing vectors. Applications using
custom adapters are responsible for permission and adapter correctness.

External signer integration, delegation authority testing and verified protocol
compatibility remain future work. Unit and race tests validate local behavior;
they are not an exchange security audit or financial risk guarantee.

# Roadmap

## Prototype — delivered offline

Evidence audit and benchmark; architecture; generic HTTP and exact decimals;
proposed typed market/account reads with synthetic fixtures; bounded local
replay and book recovery; CLI/examples; signing/trading/agent design and grant
package. Live adapters, signed actions and real WebSockets remain pending.
See [prototype status](docs/PROTOTYPE_STATUS.md) and [README status](README.md).

## Grant-funded v1

1. Obtain express written automated-access permission and engineering-approved
   contracts, environments, fixtures and vectors; ship first verified reads.
2. Implement canonical signing and supported trading on a permissioned testnet.
3. Implement documented live streams and recovery without invented guarantees.
4. Implement confirmed least-privilege delegation, expiry/revocation/rotation.
5. Publish stress/race/leak/latency results, confirmed quotas and telemetry.
6. Release with real repository/module owner, compatibility/migration policy,
   keyless read examples and testnet bot/market-maker starters.

This is a conditional roadmap. v1.0 requires accepted protocol coverage and
release gates, not simply completion of an estimated calendar. Planned pre-v1
breaking API changes should use minor versions and changelog entries; after
v1, follow semantic versioning and document supported API versions.

## Future possibilities — not implemented or funded by Full scope

A market-making framework, indexer, FIX gateway, MCP server, AI-agent tools,
Rust interoperability, metrics exporters or strategy framework could follow
validated demand. They should remain separate projects or adapters rather
than expand the core SDK prematurely. No performance, profitability or adoption
promise is made.

# Contributing

This is a community prototype. Build with Go 1.23+; local validation
uses Go 1.26.5 and CI is configured for the declared minimum plus stable.
There are no external module dependencies.

Run gofmt, go test ./..., go vet ./... and go test -race ./.... HTTP tests require
loopback listening but no external network. Run the CLI demo and check local
documentation links with python3 scripts/check_docs.py. Use small reviewable
commits with clear validation and limitations. Race-test new concurrency work.

Protocol contributions need an evidence record: exact source, access date,
classification, confidence and whether engineering confirmed the contract.
Do not equate an observed app response or another venue's protocol with an
official TxFlow guarantee. Live integrations require express written automated-
access permission and a confirmed spec. Never introduce signed production probes,
implicit credentials, unsafe retries or fabricated test vectors. Fixtures need
provenance and should be synthetic unless captured under authorized access.

Keep the core API small. Adapter wire models should map to domain types and
preserve exact decimal values and meaningful missing fields. Context cancellation,
resource bounds and failure recovery must be reviewable. Unsupported live
features belong in research/design documentation. New dependencies require a
clear benefit and license review; offline tests must remain reproducible.

MIT was chosen over Apache-2.0 for simplicity and broad reuse of original code.
Apache's explicit patent provisions could be reconsidered for future organizational
contributions, without relicensing existing contributions casually. No proprietary
TxFlow code/assets are included and the name does not imply endorsement.
The repository is maintained at [0xrariux/txflow-go](https://github.com/0xrariux/txflow-go).
Report security issues through the private route in [SECURITY.md](SECURITY.md).

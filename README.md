# txflow-go

[![Go checks](https://github.com/0xrariux/txflow-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/0xrariux/txflow-go/actions/workflows/ci.yml)

A community Go SDK prototype and grant proposal for TxFlow developer infrastructure.

**Community project. Not affiliated with or endorsed by TxFlow unless that status changes.**

## What is txflow-go?

Typed read services, precise financial domain models, tested HTTP infrastructure,
and an offline order-book demonstration, designed for a future authorized TxFlow
adapter. No live TxFlow endpoint, signer or trading implementation is shipped.

## Why?

Go backends can benefit from shared transport, precision, query and recovery
infrastructure. Reviewed official docs describe TIP components but the
[Platform API reference](https://docs.txflow.com/api-reference/readme) is
Coming Soon. This repository proposes a native Go contribution; the research
does not prove there is no existing or planned official Go SDK.

## Status

Prototype **v0.1.0-prototype**. The [terms](https://docs.txflow.com/terms/terms-of-use)
restrict automated access; written permission and engineering alignment are
prerequisites for live integration. Everything demonstrated below is explicitly
synthetic and offline. [Audit](docs/TXFLOW_DEVELOPER_SURFACE.md) ·
[Evidence](docs/SOURCES.md) · [Architecture](docs/ARCHITECTURE.md).

| Feature | Status |
|---|---|
| Generic HTTP, deadlines, response bounds, typed HTTP errors | Implemented and tested locally |
| Opt-in safe-read retries, Retry-After, interval limiter, request hooks | Implemented; no TxFlow quota assumptions |
| Decimal values and typed market/account services | Implemented proposed domain API |
| Market metadata, books, trades, ticker/funding/OI | Synthetic fixtures only; live queries pending |
| Account state, orders, fills, funding history | Synthetic fixtures only; live queries pending |
| Local order book: top-N, mid/spread, gaps and recovery | Implemented synthetic sequence policy |
| Streaming interfaces and bounded replay | Implemented offline; experimental future interface |
| Live WebSocket, reconnect/heartbeat, user streams | Pending TxFlow specification and permission |
| Signing, orders and agent wallets | Proposed / grant scope; no live submissions |
| Production compatibility and hardening | Grant scope; no production claim |

## Quick start

Go 1.23+ is required. There are no external dependencies.

Install the prototype CLI:

```sh
go install github.com/0xrariux/txflow-go/cmd/txflowctl@v0.1.0-prototype
txflowctl --demo markets
```

Or clone the repository to run the examples and tests:

```sh
git clone https://github.com/0xrariux/txflow-go.git
cd txflow-go
```

From this checkout:

```sh
go test ./...
go run ./cmd/txflowctl --demo markets
go run ./cmd/txflowctl --demo book BTC-USDC
go run ./examples/local-orderbook
```

Explicit synthetic SDK usage:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/0xrariux/txflow-go/demo"
)

func main() {
    client, err := demo.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    markets, err := client.Markets.List(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("SYNTHETIC:", markets)
}
```

To use the SDK in another Go module:

```sh
go get github.com/0xrariux/txflow-go@v0.1.0-prototype
```

`txflow.NewClient()` has no network default or implicit demo fallback; read
methods return `ErrPendingConfirmation` before I/O. Injected readers are an
extension seam, not verified TxFlow support. Custom HTTP/base/WS URLs prepare
configuration but do not activate exchange endpoints.

## Demonstrations

```sh
go run ./cmd/txflowctl --demo trades BTC-USDC
go run ./cmd/txflowctl --demo ticker BTC-USDC
go run ./cmd/txflowctl --demo account 0x0000000000000000000000000000000000000000
go run ./cmd/txflowctl --demo stream book BTC-USDC
go run ./examples/list-markets
go run ./examples/orderbook
go run ./examples/account-state
go run ./examples/websocket
```

`examples/websocket` exercises a local streaming interface and opens no socket.
Replay prints duplicate rejection, gap invalidation and recovery. Its JSON,
timestamps, values and sequences are authored fixtures, not captured exchange
data. The demo helper explicitly selects contiguous sequences; TxFlow's rules
remain unknown.

## Validation and contribution

```sh
go vet ./...
go test -race ./...
go build ./cmd/txflowctl
```

Unit tests use local httptest servers and synthetic fixtures; no external network
or private keys are needed. [Contributing](CONTRIBUTING.md) · [Security](SECURITY.md).
No live integration tests are included because the access/schema gate is unresolved.

## Roadmap and grant

[Roadmap](ROADMAP.md) · [Grant proposal](GRANT_PROPOSAL.md) ·
[Budget alternatives](docs/GRANT_BUDGET.md) ·
[Questions for TxFlow](docs/QUESTIONS_FOR_TXFLOW.md) ·
[Proposal overview](proposal/ONE_PAGER.md) · [Prototype status](docs/PROTOTYPE_STATUS.md).

MIT license is selected for broad reuse of this independently authored code;
see [LICENSE](LICENSE). No TxFlow code, brand assets or signing implementation
has been copied into this project.

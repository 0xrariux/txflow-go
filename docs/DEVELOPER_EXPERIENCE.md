# Developer experience demonstration

Today the demo proves a typed application interface while the real protocol
remains pending. No exchange endpoint is invented for the “without SDK” case.

## Without a shared integration layer

After receiving an approved contract, each developer would otherwise implement
HTTP deadlines/throttling/retries, payload construction, precise parsing,
error mapping, subscriptions, reconnection, snapshot continuity and signing
serialization. Bugs in any of those layers can propagate into bot state.

## With the current prototype

```go
client, err := demo.NewClient() // Explicit synthetic data; no network.
if err != nil { return err }
book, err := client.Markets.OrderBook(ctx, types.OrderBookRequest{
    Symbol: "BTC-USDC", Depth: 20,
})
```

The runnable version is in examples/orderbook. The local-orderbook example
adds exact mid/spread calculations, duplicate handling and recovery in a few
lines via demo.RunBook. Domain callers receive decimal values and typed errors.

## Intended funded experience

The same market/account service shape can be backed by a verified adapter.
Exact constructor/network defaults and wire mappings must be reviewed with
TxFlow. The future streaming service should manage connection/recovery;
signing should hide confirmed canonicalization. Those benefits are proposed
rather than claimed as shipped. This separation makes the pitch reproducible
without creating an unsupported trading integration.

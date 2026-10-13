# Proposed trading API

Order types, side, time-in-force, reduce-only, price/quantity, order ID and
client ID are domain design only. They have no JSON/wire mapping or endpoint.
Every Client method returns ErrPendingConfirmation locally. Context errors
are honored. No method creates an account, signs, submits, cancels or moves funds.

PlaceOrder, CancelOrder, CancelAll and ModifyOrder establish a reviewable target
interface. Batching, leverage/margin, transfers and withdrawals need separate
official action confirmation. Market-order semantics and post-only behavior
are not inferred from another exchange. See the signing research and engineering
questions before implementing any authenticated adapter.

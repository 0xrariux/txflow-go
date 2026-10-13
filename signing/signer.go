// Package signing defines an unimplemented future signer boundary.
// No payload canonicalization, EIP-712 domain, hashing or key handling exists.
package signing

import "context"

// Payload must be created from an approved protocol specification and vectors.
// Data contains signing material, never a private key or seed phrase.
type Payload struct {
	Format      string
	SpecVersion string
	Data        []byte
}
type Signature struct {
	Format string
	Bytes  []byte
}

// Signer implementations must not log secrets and must honor cancellation.
// Address is proposed as a public hex address; its network semantics await
// protocol alignment. Sign is not called by any prototype trading operation.
type Signer interface {
	Address() string
	Sign(context.Context, Payload) (Signature, error)
}

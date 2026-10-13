package main

import (
	"bytes"
	"context"
	"errors"
	txflow "github.com/0xrariux/txflow-go"
	"github.com/0xrariux/txflow-go/demo"
	"strings"
	"testing"
)

func TestCLICommandsAndNoImplicitDemo(t *testing.T) {
	for _, args := range [][]string{{"--demo", "markets"}, {"--demo", "assets"}, {"--demo", "book", "BTC-USDC"}, {"--demo", "trades", "BTC-USDC"}, {"--demo", "ticker", "BTC-USDC"}, {"--demo", "account", demo.Address}, {"--demo", "stream", "book", "BTC-USDC"}} {
		var out, diag bytes.Buffer
		if err := run(context.Background(), args, &out, &diag); err != nil {
			t.Fatal(args, err)
		}
		if out.Len() == 0 || !strings.Contains(diag.String(), "SYNTHETIC") {
			t.Fatal(args, "missing output/source label")
		}
	}
	var out, diag bytes.Buffer
	if err := run(context.Background(), []string{"markets"}, &out, &diag); !errors.Is(err, txflow.ErrPendingConfirmation) || out.Len() != 0 {
		t.Fatal("default produced unverified data", err)
	}
	for _, args := range [][]string{{}, {"--demo", "book"}, {"--demo", "stream", "orders", "BTC-USDC"}, {"--demo", "trade"}, {"--demo", "markets", "extra"}} {
		if err := run(context.Background(), args, &out, &diag); !errors.Is(err, txflow.ErrInvalidRequest) {
			t.Fatal(args, err)
		}
	}
}
func TestReplayDemonstratesRecovery(t *testing.T) {
	var out, diag bytes.Buffer
	if err := run(context.Background(), []string{"--demo", "stream", "book", "BTC-USDC"}, &out, &diag); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"sequence=100", "duplicate/stale", "sequence gap", "sequence=103", "sequence=104", "mid=70000.1"} {
		if !strings.Contains(out.String(), s) {
			t.Fatal("missing", s, out.String())
		}
	}
}

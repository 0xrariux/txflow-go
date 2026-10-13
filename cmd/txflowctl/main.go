// txflowctl is an offline SDK demonstrator; live endpoints await confirmation.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	txflow "github.com/0xrariux/txflow-go"
	"github.com/0xrariux/txflow-go/demo"
	"github.com/0xrariux/txflow-go/types"
	"io"
	"os"
	"os/signal"
)

const usage = "usage: txflowctl --demo {markets|assets|book SYMBOL|trades SYMBOL|ticker SYMBOL|account ADDRESS|stream book BTC-USDC}\nWithout --demo, live operations return pending-confirmation errors.\n"

func run(ctx context.Context, args []string, out, diagnostics io.Writer) error {
	flags := flag.NewFlagSet("txflowctl", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	offline := flags.Bool("demo", false, "use authored synthetic fixtures; no network")
	flags.Usage = func() { fmt.Fprint(diagnostics, usage) }
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	args = flags.Args()
	if len(args) == 0 {
		flags.Usage()
		return txflow.ErrInvalidRequest
	}
	command := args[0]
	arity := map[string]int{"markets": 1, "assets": 1, "book": 2, "trades": 2, "ticker": 2, "account": 2, "stream": 3}
	n, ok := arity[command]
	if !ok || len(args) != n {
		flags.Usage()
		return txflow.ErrInvalidRequest
	}
	if command == "stream" && (args[1] != "book" || args[2] != "BTC-USDC") {
		flags.Usage()
		return txflow.ErrInvalidRequest
	}
	var c *txflow.Client
	var err error
	if *offline {
		c, err = demo.NewClient()
		fmt.Fprintln(diagnostics, "SYNTHETIC OFFLINE DEMO; data is not from TxFlow")
	} else {
		c, err = txflow.NewClient()
	}
	if err != nil {
		return err
	}
	var result any
	switch command {
	case "markets":
		result, err = c.Markets.List(ctx)
	case "assets":
		result, err = c.Markets.Assets(ctx)
	case "book":
		result, err = c.Markets.OrderBook(ctx, types.OrderBookRequest{Symbol: args[1], Depth: 20})
	case "trades":
		result, err = c.Markets.Trades(ctx, args[1])
	case "ticker":
		result, err = c.Markets.Ticker(ctx, args[1])
	case "account":
		result, err = c.Account.State(ctx, args[1])
	case "stream":
		if !*offline {
			_, err = c.Stream.Connect(ctx)
			return err
		}
		return demo.RunBook(ctx, out)
	}
	if err != nil {
		return err
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

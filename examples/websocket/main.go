// Exercises the future streaming interface with LOCAL REPLAY.
// It does not open a WebSocket or claim TxFlow protocol compatibility.
package main

import (
	"context"
	"fmt"
	"github.com/0xrariux/txflow-go/demo"
	"log"
)

func main() {
	ctx := context.Background()
	c, err := demo.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	session, err := c.Stream.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	sub, err := session.SubscribeOrderBook(ctx, "BTC-USDC")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()
	fmt.Println("SYNTHETIC OFFLINE STREAM; no WebSocket connection")
	for e := range sub.Updates() {
		if e.Snapshot != nil {
			fmt.Println("snapshot", *e.Snapshot.Sequence)
		} else {
			fmt.Println("delta", e.Delta.Sequence)
		}
	}
	if err := sub.Err(); err != nil {
		log.Fatal(err)
	}
}

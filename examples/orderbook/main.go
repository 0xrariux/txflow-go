package main

import (
	"context"
	"fmt"
	"github.com/0xrariux/txflow-go/demo"
	"github.com/0xrariux/txflow-go/types"
	"log"
)

func main() {
	client, err := demo.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	book, err := client.Markets.OrderBook(context.Background(), types.OrderBookRequest{Symbol: "BTC-USDC", Depth: 2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SYNTHETIC OFFLINE DEMO", book.Symbol)
	fmt.Println("bids", book.Bids, "asks", book.Asks)
}

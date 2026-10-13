package main

import (
	"context"
	"fmt"
	"github.com/0xrariux/txflow-go/demo"
	"log"
)

func main() {
	c, err := demo.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	state, err := c.Account.State(context.Background(), demo.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SYNTHETIC OFFLINE DEMO")
	fmt.Println("equity", state.Equity, "balances", state.Balances)
}

package main

import (
	"context"
	"fmt"
	"github.com/0xrariux/txflow-go/demo"
	"log"
)

func main() {
	client, err := demo.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := client.Markets.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SYNTHETIC OFFLINE DEMO")
	for _, m := range rows {
		fmt.Println(m.Symbol, m.Kind, m.PriceIncrement)
	}
}

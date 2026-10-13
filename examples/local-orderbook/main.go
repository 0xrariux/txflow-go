package main

import (
	"context"
	"github.com/0xrariux/txflow-go/demo"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := demo.RunBook(ctx, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

package demo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/0xrariux/txflow-go/stream"
	"github.com/0xrariux/txflow-go/types"
)

//go:embed testdata/book-events.json
var bookEvents []byte

func newConnector() (*stream.ReplayConnector, error) {
	var events []types.BookEvent
	if err := json.Unmarshal(bookEvents, &events); err != nil {
		return nil, fmt.Errorf("demo: replay fixture: %w", err)
	}
	return &stream.ReplayConnector{Symbol: "BTC-USDC", Events: events, Buffer: 1}, nil
}

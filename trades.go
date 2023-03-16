package stream

import "github.com/phoebetronic/trades/typ/trades"

type Trades struct {
	Src map[string]*trades.Trades
	Dst map[string]*trades.Trades
}

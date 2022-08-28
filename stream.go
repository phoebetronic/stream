package stream

import "github.com/phoebetron/trades/typ/trades"

func Stream(str map[string]chan *trades.Trades) chan map[string]*trades.Trades {
	var tra chan map[string]*trades.Trades
	{
		tra = make(chan map[string]*trades.Trades, 1)
	}

	go func() {
		for {
			m := map[string]*trades.Trades{}
			for k, v := range str {
				m[k] = <-v
			}
			tra <- m
		}
	}()

	return tra
}

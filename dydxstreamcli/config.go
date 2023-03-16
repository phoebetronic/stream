package dydxstreamcli

import "github.com/phoebetronic/trades/typ/market"

type Config struct {
	Mar market.Market
}

func (c Config) Verify() {
	if c.Mar == nil {
		panic("Config.Mar must not be empty")
	}
}

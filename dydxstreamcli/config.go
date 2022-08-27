package dydxstreamcli

import "github.com/phoebetron/trades/typ/market"

type Config struct {
	Mar *market.Market
}

func (c Config) Verify() {
	if c.Mar == nil {
		panic("Config.Market must not be empty")
	}
}

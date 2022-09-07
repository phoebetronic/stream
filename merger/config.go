package merger

import (
	"github.com/phoebetron/stream/direct"
	"github.com/phoebetron/stream/volume"
	"github.com/phoebetron/trades/typ/trades"
)

type Config struct {
	Src map[string]chan *trades.Trades
	Dst map[string]chan *trades.Trades
	Dir direct.Config
	Vol volume.Config
}

func (c Config) Verify() {
	if len(c.Src) == 0 {
		panic("Config.Src must not be empty")
	}
	if len(c.Dst) == 0 {
		panic("Config.Dst must not be empty")
	}
	if len(c.Dst) > 1 {
		panic("Config.Dst must not contain more than one destination")
	}
}

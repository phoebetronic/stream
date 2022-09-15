package stream

import (
	"github.com/go-numb/go-dydx/public"
	"github.com/phoebetron/trades/typ/trades"
)

type Stream interface {
	Orders() chan public.OrderbookResponse
	// Trades returns a channel that streams Trades buffers with a buffer length
	// provided by the underlying market. Using a buffer length of e.g. 1 second
	// would buffer all trades that happened within the same second and send all
	// of those trades wrapped in a trades container over the returned channel
	// once the full second concluded. A trades wrapper will be send over the
	// returned channel regardless the amount of trades streamed by the
	// underlying exchange. In any case the trades wrapper will always contain
	// at least a single trade, which is the last known trade pointing to the
	// last known price.
	Trades() chan *trades.Trades
}

type Merger interface {
	Direct(key ...string) Direct
	Trades() chan Trades
	Volume(key ...string) Volume
}

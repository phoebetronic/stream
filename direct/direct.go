package direct

import (
	"time"

	"github.com/phoebetron/trades/typ/trades"
)

type Direct struct {
	dur time.Duration
	his int
	res float32
	tra []*trades.Trades
}

func New(con Config) *Direct {
	{
		con = con.Ensure()
	}

	{
		con.Verify()
	}

	return &Direct{
		dur: con.Dur,
		his: con.His,
		tra: []*trades.Trades{{}},
	}
}

func (d *Direct) Sample(tra *trades.Trades) {
	if !d.tra[len(d.tra)-1].ST.AsTime().Add(d.dur).After(tra.ST.AsTime()) {
		d.tra = append(d.tra, tra)
	}

	if len(d.tra) > d.his {
		{
			copy(d.tra[0:], d.tra[1:])
			d.tra[len(d.tra)-1] = nil
			d.tra = d.tra[:len(d.tra)-1]
		}

		{
			d.res = 0
		}

		for i := 1; i < d.his; i++ {
			d.res += res(d.tra[i-1].PR().Avg(), d.tra[i].PR().Avg())
		}
	}
}

func (d *Direct) Result() float32 {
	return d.res
}

func res(a float32, b float32) float32 {
	if a > b {
		return -1
	}

	if a < b {
		return +1
	}

	return 0
}

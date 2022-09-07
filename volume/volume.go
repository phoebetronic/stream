package volume

import (
	"github.com/phoebetron/stream"
	"github.com/phoebetron/trades/typ/trades"
)

type Volume struct {
	his int
	res stream.Volume
	vol []stream.Volume
}

func New(con Config) *Volume {
	{
		con = con.Ensure()
	}

	{
		con.Verify()
	}

	return &Volume{
		his: con.His,
	}
}

func (v *Volume) Sample(tra *trades.Trades) {
	{
		v.vol = append(v.vol, stream.Volume{Lon: tra.LO().Sum(), Sho: tra.SH().Sum()})
	}

	if len(v.vol) > v.his {
		{
			copy(v.vol[0:], v.vol[1:])
			v.vol[len(v.vol)-1] = stream.Volume{}
			v.vol = v.vol[:len(v.vol)-1]
		}

		var lon float32
		var sho float32

		for i := range v.vol {
			lon += v.vol[i].Lon
			sho += v.vol[i].Sho
		}

		{
			v.res.Lon = lon / float32(v.his)
			v.res.Sho = sho / float32(v.his)
		}
	}
}

func (v *Volume) Result() stream.Volume {
	return v.res
}

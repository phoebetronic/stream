package volume

import (
	"github.com/phoebetron/stream"
	"github.com/phoebetron/trades/typ/trades"
)

type Volume struct {
	res stream.Volume
	sli *Slicer
}

func New(con Config) *Volume {
	{
		con = con.Ensure()
	}

	{
		con.Verify()
	}

	return &Volume{
		sli: &Slicer{
			his: con.His,
		},
	}
}

func (v *Volume) Sample(tra *trades.Trades) {
	{
		var lon float32
		{
			lon = tra.LO().Sum()
		}

		var sho float32
		{
			sho = tra.SH().Sum()
		}

		var net float32
		{
			net = lon - sho
		}

		v.sli.Add(stream.Volume{
			Net: net,
			Lon: lon,
			Sho: sho,
		})
	}

	{
		var net float32
		var lon float32
		var sho float32

		for _, f := range v.sli.Lis() {
			net += f.Net
			lon += f.Lon
			sho += f.Sho
		}

		{
			v.res.Net = net
			v.res.Lon = lon
			v.res.Sho = sho
		}
	}
}

func (v *Volume) Result() stream.Volume {
	return v.res
}

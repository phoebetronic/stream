package merger

import (
	"github.com/phoebetron/stream"
	"github.com/phoebetron/stream/direct"
	"github.com/phoebetron/stream/volume"
	"github.com/phoebetron/trades/typ/trades"
)

type Merger struct {
	src map[string]chan *trades.Trades
	dst map[string]chan *trades.Trades

	dir map[string]*direct.Direct
	vol map[string]*volume.Volume
}

func New(con Config) *Merger {
	{
		con.Verify()
	}

	dir := map[string]*direct.Direct{}
	for k := range con.Src {
		dir[k] = direct.New(con.Dir)
	}

	vol := map[string]*volume.Volume{}
	for k := range con.Dst {
		vol[k] = volume.New(con.Vol)
	}

	return &Merger{
		src: con.Src,
		dst: con.Dst,
		dir: dir,
		vol: vol,
	}
}

func (m *Merger) Direct() stream.Direct {
	var sum float32

	for _, v := range m.dir {
		sum += v.Result()
	}

	return stream.Direct{Dir: sum / float32(len(m.dir))}
}

func (m *Merger) Trades() chan stream.Trades {
	var tra chan stream.Trades
	{
		tra = make(chan stream.Trades, 1)
	}

	go func() {
		for {
			t := stream.Trades{
				Src: map[string]*trades.Trades{},
				Dst: map[string]*trades.Trades{},
			}

			for k, v := range m.src {
				var p *trades.Trades
				{
					p = <-v
				}

				{
					t.Src[k] = p
				}

				{
					m.dir[k].Sample(p)
				}
			}

			for k, v := range m.dst {
				var p *trades.Trades
				{
					p = <-v
				}

				{
					t.Dst[k] = p
				}

				{
					m.vol[k].Sample(p)
				}
			}

			{
				tra <- t
			}
		}
	}()

	return tra
}

func (m *Merger) Volume() stream.Volume {
	var lon float32
	var sho float32

	for _, v := range m.vol {
		lon += v.Result().Lon
		sho += v.Result().Sho
	}

	return stream.Volume{Lon: lon / float32(len(m.vol)), Sho: sho / float32(len(m.vol))}
}

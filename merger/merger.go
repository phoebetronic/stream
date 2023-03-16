package merger

import (
	"github.com/phoebetronic/stream"
	"github.com/phoebetronic/stream/direct"
	"github.com/phoebetronic/stream/volume"
	"github.com/phoebetronic/trades/typ/trades"
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
	vol := map[string]*volume.Volume{}

	for k := range con.Src {
		dir[k] = direct.New(con.Dir)
		vol[k] = volume.New(con.Vol)
	}

	for k := range con.Dst {
		dir[k] = direct.New(con.Dir)
		vol[k] = volume.New(con.Vol)
	}

	return &Merger{
		src: con.Src,
		dst: con.Dst,
		dir: dir,
		vol: vol,
	}
}

func (m *Merger) Direct(key ...string) stream.Direct {
	var sum float32

	for k, v := range m.dir {
		if !contains(key, k) {
			continue
		}

		{
			sum += v.Result()
		}
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
					m.vol[k].Sample(p)
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
					m.dir[k].Sample(p)
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

func (m *Merger) Volume(key ...string) stream.Volume {
	var net float32
	var lon float32
	var sho float32

	for k, v := range m.vol {
		if !contains(key, k) {
			continue
		}

		{
			net += v.Result().Net
			lon += v.Result().Lon
			sho += v.Result().Sho
		}
	}

	return stream.Volume{
		Net: net / float32(len(m.vol)),
		Lon: lon / float32(len(m.vol)),
		Sho: sho / float32(len(m.vol)),
	}
}

func contains(lis []string, ele string) bool {
	for _, e := range lis {
		if e == ele {
			return true
		}
	}

	return false
}

package dydxstreamcli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-numb/go-dydx/public"
	"github.com/go-numb/go-dydx/realtime"
	"github.com/phoebetron/trades/typ/buffer"
	"github.com/phoebetron/trades/typ/market"
	"github.com/phoebetron/trades/typ/trades"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Stream struct {
	buf buffer.Interface
	mar market.Interface
}

func New(con Config) *Stream {
	{
		con.Verify()
	}

	var buf buffer.Interface
	{
		buf = buffer.New(buffer.Config{
			Mar: con.Mar,
		})
	}

	return &Stream{
		buf: buf,
		mar: con.Mar,
	}
}

func (s *Stream) Trades() chan *trades.Trades {
	res := make(chan realtime.Response)

	go func() {
		err := realtime.Connect(
			context.Background(),
			res,
			[]string{realtime.TRADES},
			[]string{fmt.Sprintf("%s-USD", strings.ToUpper(s.mar.Ass()))},
			nil,
			nil,
		)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			var cur time.Time
			{
				cur = time.Now().UTC()
			}

			var dur time.Duration
			{
				dur = cur.Truncate(s.mar.Dur()).Add(s.mar.Dur()).Sub(cur)
			}

			{
				time.Sleep(dur)
			}

			{
				s.buf.Finish(time.Now().UTC())
			}
		}
	}()

	go func() {
		for re := range res {
			switch re.Channel {
			case realtime.TRADES:
				s.trades(re.Trades.Trades)
			case realtime.ERROR:
				panic(re.Results)
			case realtime.UNDEFINED:
				panic(re.Results)
			}
		}
	}()

	return s.buf.Trades()
}

func (s *Stream) trades(raw []public.Trade) {
	for _, r := range raw {
		t := &trades.Trade{}
		{
			t.PR = musf32(r.Price)
			t.TS = timestamppb.New(r.CreatedAt)
		}

		if strings.ToLower(r.Side) == "buy" {
			t.LO = musf32(r.Size)
		}

		if strings.ToLower(r.Side) == "sell" {
			t.SH = musf32(r.Size)
		}

		{
			s.buf.Buffer(t)
		}
	}
}

func musf32(s string) float32 {
	f, e := strconv.ParseFloat(s, 32)
	if e != nil {
		panic(e)
	}

	return float32(f)
}

package ftxstreamcli

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/go-numb/go-ftx/realtime"
	"github.com/go-numb/go-ftx/rest/public/markets"
	"github.com/phoebetron/trades/typ/buffer"
	"github.com/phoebetron/trades/typ/market"
	"github.com/phoebetron/trades/typ/trades"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Stream struct {
	buf buffer.Buffer
	mar market.Market
}

func New(con Config) *Stream {
	{
		con.Verify()
	}

	var buf buffer.Buffer
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
			[]string{"trades"},
			[]string{fmt.Sprintf("%s-PERP", strings.ToUpper(s.mar.Ass()))},
			log.New(ioutil.Discard, "", 0),
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
			switch re.Type {
			case realtime.TRADES:
				s.trades(re.Trades)
			case realtime.ERROR:
				panic(re.Results)
			case realtime.UNDEFINED:
				panic(re.Results)
			}
		}
	}()

	return s.buf.Trades()
}

func (s *Stream) trades(raw []markets.Trade) {
	for _, r := range raw {
		t := &trades.Trade{}
		{
			t.PR = float32(r.Price)
			t.TS = timestamppb.New(r.Time)
		}

		if strings.ToLower(r.Side) == "buy" {
			t.LO = float32(r.Size)
		}

		if strings.ToLower(r.Side) == "sell" {
			t.SH = float32(r.Size)
		}

		{
			s.buf.Buffer(t)
		}
	}
}

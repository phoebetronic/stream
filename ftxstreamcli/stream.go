package ftxstreamcli

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/go-numb/go-dydx/public"
	"github.com/go-numb/go-ftx/realtime"
	"github.com/go-numb/go-ftx/rest/public/markets"
	"github.com/phoebetron/ftxapi/client"
	"github.com/phoebetron/ftxapi/client/public/trade"
	"github.com/phoebetron/trades/typ/buffer"
	"github.com/phoebetron/trades/typ/market"
	"github.com/phoebetron/trades/typ/trades"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Stream struct {
	buf buffer.Buffer
	cli *client.Client
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

	var cli *client.Client
	{
		cli = client.New(client.Config{})
	}

	return &Stream{
		buf: buf,
		cli: cli,
		mar: con.Mar,
	}
}

func (s *Stream) Orders() chan public.OrderbookResponse {
	panic("Stream.Orders is not implemented for FTX")
}

func (s *Stream) Trades() chan *trades.Trades {

	go func() {
		var err error

		var req trade.ListRequest
		{
			req = trade.ListRequest{
				ProductCode: fmt.Sprintf("%s-PERP", strings.ToUpper(s.mar.Ass())),
				Limit:       1,
			}
		}

		var res trade.ListResponse
		{
			res, err = s.cli.Pub.Tra.List(req)
			if err != nil {
				panic(err)
			}
		}

		s.buf.Latest(mustra(markets.Trade{
			Side:  res.Result[0].Side,
			Size:  res.Result[0].Size,
			Price: res.Result[0].Price,
			Time:  res.Result[0].Time,
		}))
	}()

	go func() {
		for {
			clo := make(chan struct{})
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
				for re := range res {
					switch re.Type {
					case realtime.TRADES:
						for _, r := range re.Trades {
							s.buf.Buffer(mustra(r))
						}
					case realtime.ERROR:
						close(clo)
						return
					case realtime.UNDEFINED:
						close(clo)
						return
					}
				}
			}()

			{
				<-clo
			}

			{
				fmt.Printf("restarting ftx websocket\n")
			}
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

	return s.buf.Trades()
}

func mustra(r markets.Trade) *trades.Trade {
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

	return t
}

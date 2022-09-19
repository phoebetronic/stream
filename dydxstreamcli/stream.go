package dydxstreamcli

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-numb/go-dydx/public"
	"github.com/go-numb/go-dydx/realtime"
	"github.com/phoebetron/dydxv3/client"
	"github.com/phoebetron/dydxv3/client/public/trade"
	"github.com/phoebetron/trades/typ/market"
	"github.com/phoebetron/trades/typ/trades"
	"github.com/phoebetron/trades/typ/trades/buffer"
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
	ord := make(chan public.OrderbookResponse)

	go func() {
		for {
			clo := make(chan struct{})
			res := make(chan realtime.Response)

			go func() {
				err := realtime.Connect(
					context.Background(),
					res,
					[]string{realtime.ORDERBOOK},
					[]string{fmt.Sprintf("%s-USD", strings.ToUpper(s.mar.Ass()))},
					nil,
					log.New(ioutil.Discard, "", 0),
				)
				if err != nil {
					panic(err)
				}
			}()

			go func() {
				for re := range res {
					switch re.Channel {
					case realtime.ORDERBOOK:
						ord <- re.Orderbook
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
				fmt.Printf("restarting dydx websocket\n")
			}
		}
	}()

	return ord
}

func (s *Stream) Trades() chan *trades.Trades {
	go func() {
		var err error

		var req trade.ListRequest
		{
			req = trade.ListRequest{
				Market: fmt.Sprintf("%s-USD", strings.ToUpper(s.mar.Ass())),
				Limit:  1,
			}
		}

		var res trade.ListResponse
		{
			res, err = s.cli.Pub.Tra.List(req)
			if err != nil {
				panic(err)
			}
		}

		s.buf.Latest(mustra(public.Trade{
			Side:      res.Trades[0].Side,
			Size:      res.Trades[0].Size,
			Price:     res.Trades[0].Price,
			CreatedAt: res.Trades[0].CreatedAt,
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
					[]string{realtime.TRADES},
					[]string{fmt.Sprintf("%s-USD", strings.ToUpper(s.mar.Ass()))},
					nil,
					log.New(ioutil.Discard, "", 0),
				)
				if err != nil {
					panic(err)
				}
			}()

			go func() {
				for re := range res {
					switch re.Channel {
					case realtime.TRADES:
						for _, r := range re.Trades.Trades {
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
				fmt.Printf("restarting dydx websocket\n")
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

func musf32(s string) float32 {
	f, e := strconv.ParseFloat(s, 32)
	if e != nil {
		panic(e)
	}

	return float32(f)
}

func mustra(r public.Trade) *trades.Trade {
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

	return t
}

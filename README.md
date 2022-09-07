# stream

```go
package main

import (
	"fmt"
	"time"

	"github.com/phoebetron/stream/dydxstreamcli"
	"github.com/phoebetron/stream/ftxstreamcli"
	"github.com/phoebetron/stream/merger"
	"github.com/phoebetron/trades/typ/market"
	"github.com/phoebetron/trades/typ/trades"
)

func main() {
	dyd := dydxstreamcli.New(dydxstreamcli.Config{
		Mar: market.New(market.Config{
			Exc: "dydx",
			Ass: "eth",
			Dur: 1 * time.Second,
		}),
	})

	ftx := ftxstreamcli.New(ftxstreamcli.Config{
		Mar: market.New(market.Config{
			Exc: "ftx",
			Ass: "eth",
			Dur: 1 * time.Second,
		}),
	})

	str := merger.New(merger.Config{
		Src: map[string]chan *trades.Trades{
			"ftx": ftx.Trades(),
		},
		Dst: map[string]chan *trades.Trades{
			"dydx": dyd.Trades(),
		},
	})

	for t := range str.Trades() {
		for _, v := range t.Src {
			fmt.Printf("src (%s/%s)\n", v.EX, v.AS)
			for _, t := range v.TR {
				fmt.Printf("    TR %f (%s)\n", t.PR, t.TS.AsTime())
			}
		}

		for _, v := range t.Dst {
			fmt.Printf("dst (%s/%s)\n", v.EX, v.AS)
			for _, t := range v.TR {
				fmt.Printf("    TR %f (%s)\n", t.PR, t.TS.AsTime())
			}
		}

		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("\n")
	}
}

```

```
$ go run main.go
src (ftx/eth)
    TR 1573.500000 (2022-09-05 13:04:57.215866 +0000 UTC)
    TR 1573.500000 (2022-09-05 13:04:57.216515 +0000 UTC)
dst (dydx/eth)
    TR 1573.400024 (2022-09-05 13:04:57 +0000 UTC)



src (ftx/eth)
    TR 1573.599976 (2022-09-05 13:04:58.050723 +0000 UTC)
dst (dydx/eth)
    TR 1573.400024 (2022-09-05 13:04:58 +0000 UTC)



src (ftx/eth)
    TR 1573.500000 (2022-09-05 13:04:59.234147 +0000 UTC)
dst (dydx/eth)
    TR 1573.599976 (2022-09-05 13:04:59.549 +0000 UTC)
    TR 1573.699951 (2022-09-05 13:04:59.549 +0000 UTC)



...
```

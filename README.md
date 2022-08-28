# stream

```go
package main

import (
	"fmt"
	"time"

	"github.com/phoebetron/stream/dydxstreamcli"
	"github.com/phoebetron/stream/ftxstreamcli"
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

	str := Stream(map[string]chan *trades.Trades{
		"dydx": dyd.Trades(),
		"ftx":  ftx.Trades(),
	})

	for m := range str {
		for k, v := range m {
			fmt.Printf("%s (%s/%s)\n", k, v.EX, v.AS)
			for _, t := range v.TR {
				fmt.Printf("    TS %s\n", t.TS.AsTime())
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
dydx (dydx/eth)
    TS 2022-08-28 21:45:43.067 +0000 UTC
    TS 2022-08-28 21:45:43.068 +0000 UTC
    TS 2022-08-28 21:45:43.097 +0000 UTC
    ...
ftx (ftx/eth)
    TS 2022-08-28 21:45:43.099344 +0000 UTC
    TS 2022-08-28 21:45:43.100027 +0000 UTC
    TS 2022-08-28 21:45:43.101325 +0000 UTC
    ...



ftx (ftx/eth)
    TS 2022-08-28 21:45:44.036106 +0000 UTC
    TS 2022-08-28 21:45:44.036584 +0000 UTC
    TS 2022-08-28 21:45:44.036584 +0000 UTC
    ...
dydx (dydx/eth)
    TS 2022-08-28 21:45:44.124 +0000 UTC
    TS 2022-08-28 21:45:44.165 +0000 UTC
    TS 2022-08-28 21:45:44.166 +0000 UTC
    ...



ftx (ftx/eth)
    TS 2022-08-28 21:45:45.028686 +0000 UTC
    TS 2022-08-28 21:45:45.030149 +0000 UTC
    TS 2022-08-28 21:45:45.030749 +0000 UTC
    ...
dydx (dydx/eth)
    TS 2022-08-28 21:45:45.134 +0000 UTC
    TS 2022-08-28 21:45:45.134 +0000 UTC
    TS 2022-08-28 21:45:45.134 +0000 UTC
    ...
```

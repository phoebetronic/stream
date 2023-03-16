# stream

```go
package main

import (
	"fmt"
	"time"

	"github.com/phoebetronic/stream/dydxstreamcli"
	"github.com/phoebetronic/stream/merger"
	"github.com/phoebetronic/trades/typ/market"
	"github.com/phoebetronic/trades/typ/trades"
)

func main() {
	dyd := dydxstreamcli.New(dydxstreamcli.Config{
		Mar: market.New(market.Config{
			Exc: "dydx",
			Ass: "eth",
			Dur: 250 * time.Millisecond,
		}),
	})

	str := merger.New(merger.Config{
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
dst (dydx/eth)
    TR 1573.400024 (2022-09-05 13:04:57 +0000 UTC)



dst (dydx/eth)
    TR 1573.400024 (2022-09-05 13:04:58 +0000 UTC)



dst (dydx/eth)
    TR 1573.599976 (2022-09-05 13:04:59.549 +0000 UTC)
    TR 1573.699951 (2022-09-05 13:04:59.549 +0000 UTC)



...
```

# stream

```go
import (
	"github.com/phoebetron/stream/dydxstreamcli"
	"github.com/phoebetron/trades/typ/market"
	"github.com/phoebetron/trades/typ/trades"
)

str := dydxstreamcli.New(dydxstreamcli.Config{
	Mar: market.New(market.Config{
		Exc: "dydx",
		Ass: "eth",
		Dur: 1 * time.Second,
	}),
})

for t := range str.Trades() {
	fmt.Printf("lo %s\n", time.Now().UTC())
	fmt.Printf("ST %s\n", t.ST.AsTime())
	fmt.Printf("EN %s\n", t.EN.AsTime())
	fmt.Printf("LA %f\n", t.LA().PR)
	for _, v := range t.TR {
		fmt.Printf("    TS %s\n", v.TS.AsTime())
	}
	fmt.Printf("\n")
}
```

```
lo 2022-08-27 14:33:26.000377 +0000 UTC
ST 2022-08-27 14:33:25 +0000 UTC
EN 2022-08-27 14:33:26 +0000 UTC
LA 1500.599976
    TS 2022-08-27 14:33:25 +0000 UTC

lo 2022-08-27 14:33:27.001229 +0000 UTC
ST 2022-08-27 14:33:26 +0000 UTC
EN 2022-08-27 14:33:27 +0000 UTC
LA 1500.500000
    TS 2022-08-27 14:33:26.142 +0000 UTC
    TS 2022-08-27 14:33:26.142 +0000 UTC
    TS 2022-08-27 14:33:26.154 +0000 UTC
    TS 2022-08-27 14:33:26.154 +0000 UTC

lo 2022-08-27 14:33:28.000342 +0000 UTC
ST 2022-08-27 14:33:27 +0000 UTC
EN 2022-08-27 14:33:28 +0000 UTC
LA 1500.500000
    TS 2022-08-27 14:33:27 +0000 UTC
```

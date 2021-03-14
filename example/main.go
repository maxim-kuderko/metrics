package main

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
	"github.com/un000/grpc-snappy"
	"github.com/valyala/fastrand"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	con, _ := strconv.Atoi(os.Getenv(`CON`))
	card, _ := strconv.Atoi(os.Getenv(`CARD`))
	c(con, card)
}
func c(c, card int) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	reporter := metrics.NewReporter(metrics.WithFlushTicker(time.Millisecond*10), metrics.WithDriver(drivers.NewGrpc(ctx, `localhost:8081`, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor(snappy.Name)))))
	wg := sync.WaitGroup{}
	wg.Add(c)
	for i := 0; i < c; i++ {
		go func() {
			defer wg.Done()
			for {
				reporter.Send(fmt.Sprintf(`aa%d`, fastrand.Uint32n(uint32(card))), 1)
			}
		}()
	}
	wg.Wait()
}

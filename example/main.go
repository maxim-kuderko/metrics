package main

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
	snappy "github.com/un000/grpc-snappy"
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
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	reporter := metrics.NewReporter(metrics.WithDriver(drivers.NewGrpc(ctx, `localhost:8081`, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor(snappy.Name))), 10000))
	wg := sync.WaitGroup{}
	wg.Add(c)
	names := make([]string, 0, card)
	for i := 0; i < card; i++ {
		names = append(names, fmt.Sprintf(`aa%d`, fastrand.Uint32n(uint32(card))))
	}
	for i := 0; i < c; i++ {
		go func() {
			j := 0
			defer wg.Done()
			for {
				reporter.Send(`metric_name`, 1, `tag`, names[j])
				if j == len(names)-1 {
					j = 0
				} else {
					j++
				}
			}
		}()
	}
	wg.Wait()
}

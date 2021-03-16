package main

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
	snappy "github.com/un000/grpc-snappy"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	con, _ := strconv.Atoi(os.Getenv(`CON`))
	card, _ := strconv.Atoi(os.Getenv(`CARD`))
	c(con, card)
}
func c(c, card int) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	reporter := metrics.NewReporter(metrics.WithDriver(func() metrics.Driver {
		return drivers.NewGrpc(ctx, `localhost:8081`, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor(snappy.Name)))
	}, c))
	/*reporter := metrics.NewReporter(metrics.WithDriver(func() metrics.Driver {
		return drivers.NewHTTP(`http://localhost:8080/send`, time.Second)
	}, c))*/
	wg := sync.WaitGroup{}
	wg.Add(c)
	names := make([]string, 0, card)
	for i := 0; i < card; i++ {
		names = append(names, fmt.Sprintf(`aa%d`, rand.Int31n(int32(card))))
	}
	for i := 0; i < c; i++ {
		go func() {
			j := 0
			defer wg.Done()
			for {
				reporter.Send(names[j], 1)
				if j > len(names) {
					j = 0
				}
			}
		}()
	}
	wg.Wait()
}

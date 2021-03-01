package main

import (
	"fmt"
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
	"github.com/valyala/fastrand"
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
	reporter := metrics.NewReporter(metrics.WithDriver(drivers.NewHTTP(`http://192.168.0.185:8080/send`, time.Second*10)))
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

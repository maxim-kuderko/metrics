package main

import (
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
	"time"
)

func main() {
	reporter := metrics.NewReporter(metrics.WithDriver(drivers.NewHTTP(`http://localhost:8080/send`, time.Second*10)), metrics.WithConcurrency(1))
	reporter.Send(`aa`, 1)
	reporter.Close()
}

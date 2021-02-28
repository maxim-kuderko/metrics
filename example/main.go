package main

import (
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
)

func main() {
	reporter := metrics.NewReporter(metrics.WithDriver(drivers.NewStdout()))
	reporter.Send(`aa`, 1)
}

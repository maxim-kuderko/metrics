package main

import (
	"fmt"
	"github.com/maxim-kuderko/metrics"
	"github.com/maxim-kuderko/metrics/drivers"
	"math/rand"
	"time"
)

func main() {
	reporter := metrics.NewReporter(metrics.WithDriver(drivers.NewHTTP(`http://localhost:8080/send`, time.Second*10)), metrics.WithConcurrency(1))
	str := RandStringRunes(20)
	for {
		reporter.Send(fmt.Sprintf(`%s%d`, str, rand.Int()%1000000), 1)
	}
	reporter.Close()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

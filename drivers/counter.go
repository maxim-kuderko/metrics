package drivers

import (
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"go.uber.org/atomic"
	"time"
)

var counter *Counter

type Counter struct {
	c *atomic.Int64
}

func (s *Counter) Send(r *proto.Metric) {
	s.c.Add(r.Values.Count)
}

func NewCounter() *Counter {
	if counter != nil {
		return counter
	}
	s := &Counter{c: atomic.NewInt64(0)}
	go func() {
		w := 1
		t := time.NewTicker(time.Second * time.Duration(w))
		for range t.C {
			fmt.Println(fmt.Sprintf("%0.2fm req/sec ", float64(s.c.Swap(0))/1000000/float64(w)))
		}
	}()
	counter = s
	return s
}

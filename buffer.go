package metrics

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
	"sync"
	"time"
)

type requestBuffer struct {
	data *proto.MetricsRequest
	idx  int

	mu     *sync.Mutex
	driver Driver
}

func newRequestBuffer(size int, driver Driver) *requestBuffer {
	vals := make([]*proto.Metric, 0, size)
	for i := 0; i < size; i++ {
		vals = append(vals, &proto.Metric{
			Values: &proto.Values{},
		})
	}
	return &requestBuffer{
		data:   &proto.MetricsRequest{Metrics: vals},
		driver: driver,
		mu:     &sync.Mutex{},
	}
}

func (rb *requestBuffer) add(name string, value float64, h uint64, tags ...string) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.data.Metrics[rb.idx].Name = name
	rb.data.Metrics[rb.idx].Tags = tags
	rb.data.Metrics[rb.idx].Values.Count = 1
	rb.data.Metrics[rb.idx].Values.Sum = value
	rb.data.Metrics[rb.idx].Values.Min = value
	rb.data.Metrics[rb.idx].Values.Max = value
	rb.data.Metrics[rb.idx].Values.First = value
	rb.data.Metrics[rb.idx].Values.Last = value
	rb.data.Metrics[rb.idx].Hash = h
	rb.data.Metrics[rb.idx].Time = time.Now().UnixNano()
	if rb.idx+1 == cap(rb.data.Metrics) {
		rb.driver.Send(rb.data)
		rb.idx = 0
		return
	}
	rb.idx++
}

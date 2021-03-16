package metrics

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
	"sync"
	"time"
)

type requestBuffer struct {
	data *proto.MetricsRequest
	pool *sync.Pool
	idx  int

	mu     *sync.Mutex
	driver Driver
}

func newRequestBuffer(size int, driver Driver) *requestBuffer {
	pool := newMetricRequestPool(size)
	return &requestBuffer{
		data:   pool.Get().(*proto.MetricsRequest),
		pool:   pool,
		driver: driver,
		mu:     &sync.Mutex{},
	}
}

func newMetricRequestPool(size int) *sync.Pool {
	return &sync.Pool{New: func() interface{} {
		vals := make([]*proto.Metric, 0, size)
		for i := 0; i < size; i++ {
			vals = append(vals, &proto.Metric{
				Values: &proto.Values{},
			})
		}
		return &proto.MetricsRequest{Metrics: vals}
	}}
}

func (rb *requestBuffer) add(name string, value float64, tags ...string) {
	h := calcHash(name, tags...)
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
		tmp := rb.data
		rb.data = rb.pool.Get().(*proto.MetricsRequest)
		go func() {
			rb.driver.Send(tmp)
			rb.pool.Put(tmp)
		}()
		rb.idx = 0
		return
	}
	rb.idx++
}

func (rb *requestBuffer) close() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.driver.Send(rb.data)
	rb.idx = 0
	return
}

package metrics

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
	"sync"
	"time"
)

func WithDriver(d Driver) Option {
	return func(r *Reporter) {
		r.driver = d
	}
}

func WithFlushTicker(duration time.Duration) Option {
	return func(r *Reporter) {
		r.ticker = time.NewTicker(duration)
	}
}

func WithConcurrency(concurrency int) Option {
	return func(r *Reporter) {
		output := make([]*sync.Mutex, concurrency)
		buffArr := make([]*proto.MetricsRequest, concurrency)
		for i := 0; i < concurrency; i++ {
			output[i] = &sync.Mutex{}
			buffArr[i] = &proto.MetricsRequest{}
		}
		r.mu = output
		r.buff = buffArr
		r.flushSemaphore = make(chan struct{}, concurrency)
	}
}

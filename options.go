package metrics

import (
	"github.com/maxim-kuderko/metrics/entities"
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
		r.bufferFlushTicker = time.NewTicker(duration)
	}
}

func WithConcurrency(concurrency int) Option {
	return func(r *Reporter) {
		output := make([]*sync.Mutex, concurrency)
		buffArr := make([]entities.Metrics, concurrency)
		for i := 0; i < concurrency; i++ {
			output[i] = &sync.Mutex{}
			buffArr[i] = entities.Metrics{}
		}
		r.mu = output
		r.buff = buffArr
		r.flushSemaphore = make(chan struct{}, concurrency)
	}
}

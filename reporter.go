package metrics

import (
	"github.com/cespare/xxhash"
	"github.com/maxim-kuderko/metrics/drivers"
	"github.com/maxim-kuderko/metrics/entities"
	"github.com/valyala/bytebufferpool"
	"runtime"
	"sync"
	"time"
)

type Reporter struct {
	driver Driver

	buff              []entities.Metrics
	bufferFlushTicker *time.Ticker

	mu             []*sync.Mutex
	idx            int
	flushSemaphore chan struct{}
	wg             sync.WaitGroup
}

type Option func(r *Reporter)

var defaultConfigs = []Option{WithDriver(drivers.NewStdout()), WithFlushTicker(time.Second), WithConcurrency(runtime.NumCPU() * 4)}

var metricsPool = sync.Pool{New: newBuff()}

func NewReporter(opt ...Option) *Reporter {
	m := &Reporter{}
	for _, o := range defaultConfigs {
		o(m)
	}
	for _, o := range opt {
		o(m)
	}
	return m
}

func (r *Reporter) flusher() {
	for {
		<-r.bufferFlushTicker.C
		for i, mu := range r.mu {
			mu.Lock()
			r.flush(i)
			mu.Unlock()
		}

	}
}

func newBuff() func() interface{} {
	return func() interface{} {
		return entities.Metrics{}
	}
}

func (r *Reporter) Send(name string, value float64, tags ...string) {
	h := calcHash(name, tags...)
	shard := h % uint64(len(r.mu))
	r.mu[shard].Lock()
	defer r.mu[shard].Unlock()
	v, ok := r.buff[shard][h]
	if !ok {
		tmp := &entities.AggregatedMetric{
			Name:   name,
			Tags:   tags,
			Values: entities.Values{},
		}
		r.buff[shard][h] = tmp
		v = tmp
	}
	v.Add(value)
}

func calcHash(name string, tags ...string) uint64 {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	b.WriteString(name)
	for _, s := range tags {
		b.WriteString(s)
	}
	return xxhash.Sum64(b.Bytes())
}
func (r *Reporter) Close() {
	for i, mu := range r.mu {
		mu.Lock()
		r.flush(i)
		mu.Unlock()
	}
	r.wg.Wait()
}

func (r *Reporter) flush(i int) {
	r.wg.Add(1)
	tmp := r.buff[i]
	r.buff[i] = metricsPool.Get().(entities.Metrics)
	r.flushSemaphore <- struct{}{}
	go func() {
		defer func() {
			tmp.Reset()
			metricsPool.Put(tmp)
			r.wg.Done()
			<-r.flushSemaphore
		}()
		r.driver.Send(tmp)
	}()
}

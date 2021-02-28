package metrics

import (
	"github.com/maxim-kuderko/metrics/drivers"
	"runtime"
	"sync"
	"time"
)

type Reporter struct {
	driver Driver

	buff                buffer
	buffSize            int
	bufferFlushTicker   *time.Ticker
	flushTickerDuration time.Duration

	mu             sync.Mutex
	flushSemaphore chan struct{}
}

type buffer struct {
	drivers.Metrics
	CreatedAt time.Time
}

type Option func(r *Reporter)

var defaultConfigs = []Option{WithDriver(drivers.NewStdout()), WithBuffer(100), WithFlushTicker(time.Second)}

var metricsPool sync.Pool

func NewReporter(opt ...Option) *Reporter {
	m := &Reporter{}
	for _, o := range defaultConfigs {
		o(m)
	}
	for _, o := range opt {
		o(m)
	}
	metricsPool = sync.Pool{New: newBuff(m.buffSize)}
	m.buff = metricsPool.Get().(buffer)
	m.flushSemaphore = make(chan struct{}, runtime.GOMAXPROCS(0))
	return m
}

func (r *Reporter) flusher() {
	for {
		<-r.bufferFlushTicker.C
		r.mu.Lock()
		r.flush()
		r.mu.Unlock()
	}
}

func newBuff(size int) func() interface{} {
	return func() interface{} {
		return buffer{
			Metrics:   make(drivers.Metrics, 0, size),
			CreatedAt: time.Now(),
		}
	}
}

func (r *Reporter) Send(name string, value float64, tags ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buff.Metrics = append(r.buff.Metrics, drivers.Metric{
		Name:  name,
		Value: value,
		Tags:  tags,
	})

	if len(r.buff.Metrics) == r.buffSize {
		r.flush()
	}
}

func (r *Reporter) flush() {
	tmp := r.buff
	r.buff = metricsPool.Get().(buffer)
	r.buff.CreatedAt = time.Now()
	r.flushSemaphore <- struct{}{}
	r.bufferFlushTicker.Reset(r.flushTickerDuration)
	go func() {
		<-r.flushSemaphore
		r.driver.Send(tmp.Metrics)
		tmp.Metrics = tmp.Metrics[:0]
		metricsPool.Put(tmp)
	}()
}

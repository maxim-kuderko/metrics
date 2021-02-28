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
	idx            int
	flushSemaphore chan struct{}
	wg             sync.WaitGroup
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
			Metrics:   make(drivers.Metrics, size),
			CreatedAt: time.Now(),
		}
	}
}

func (r *Reporter) Send(name string, value float64, tags ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m := r.buff.Metrics[r.idx]
	m.Value = value
	m.Name = name
	m.Tags = tags
	r.buff.Metrics[r.idx] = m
	r.idx++
	if r.idx >= r.buffSize-1 {
		r.flush()
	}
}
func (r *Reporter) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flush()
	r.wg.Wait()
}

func (r *Reporter) flush() {
	tmp := r.buff
	tmp.Metrics = tmp.Metrics[:r.idx]
	r.buff = metricsPool.Get().(buffer)
	r.buff.CreatedAt = time.Now()
	r.idx = 0
	r.flushSemaphore <- struct{}{}
	r.bufferFlushTicker.Reset(r.flushTickerDuration)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		<-r.flushSemaphore
		r.driver.Send(tmp.Metrics)
		metricsPool.Put(tmp)
	}()
}

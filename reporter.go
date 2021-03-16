package metrics

import (
	"github.com/cespare/xxhash"
	"github.com/valyala/bytebufferpool"
	"go.uber.org/atomic"
)

type Reporter struct {
	buff []*requestBuffer
	lb   *atomic.Int32
}

type Option func(r *Reporter)

func NewReporter(opt ...Option) *Reporter {
	m := &Reporter{
		lb: atomic.NewInt32(0),
	}
	for _, o := range opt {
		o(m)
	}
	return m
}

func (r *Reporter) Send(name string, value float64, tags ...string) {
	lb := r.lb.Inc()
	l := int32(len(r.buff))
	r.buff[lb%l].add(name, value, tags...)
	if lb > l {
		r.lb.Store(0)
	}
}

func (r *Reporter) Close() {
	for _, d := range r.buff {
		d.close()
	}
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

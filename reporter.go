package metrics

import (
	"github.com/cespare/xxhash"
	"github.com/valyala/bytebufferpool"
	"sync"
)

type Reporter struct {
	buff        *sync.Pool
	defaultTags []string
}

type Option func(r *Reporter)

func NewReporter(opt ...Option) *Reporter {
	m := &Reporter{}
	for _, o := range opt {
		o(m)
	}
	return m
}

func (r *Reporter) Send(name string, value float64, tags ...string) {
	b := r.buff.Get().(*requestBuffer)
	defer r.buff.Put(b)
	b.add(name, value, calcHash(name, tags...), r.defaultTags, tags...)
}

func (r *Reporter) Close() {

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

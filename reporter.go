package metrics

import (
	"github.com/cespare/xxhash"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/valyala/bytebufferpool"
	"time"
)

type Reporter struct {
	driver []Driver
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
	h := calcHash(name, tags...)
	m := proto.MetricPool.Get().(*proto.Metric)
	defer proto.MetricPool.Put(m)
	m.Name = name
	m.Tags = tags
	m.Values.Count = 1
	m.Values.Sum = value
	m.Values.Min = value
	m.Values.Max = value
	m.Values.First = value
	m.Values.Last = value
	m.Hash = h
	m.Time = time.Now().UnixNano()
	r.driver[h%uint64(len(r.driver))].Send(m)
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

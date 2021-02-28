package entities

import (
	"github.com/cespare/xxhash"
	"github.com/valyala/bytebufferpool"
)

type Metric struct {
	Name  string
	Value float64
	Tags  []string

	hash uint64
}

type Metrics []Metric

func (m *Metric) Hash() uint64 {
	if m.hash != 0 {
		return m.hash
	}
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	b.WriteString(m.Name)
	for _, t := range m.Tags {
		b.WriteString(t)
	}

	m.hash = xxhash.Sum64(b.Bytes())
	return m.hash
}

type AggregatedMetric struct {
}

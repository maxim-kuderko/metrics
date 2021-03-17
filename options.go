package metrics

import "sync"

func WithDriver(d Driver, bulkSize int) Option {
	return func(r *Reporter) {
		r.buff = &sync.Pool{New: func() interface{} { return newRequestBuffer(bulkSize, d) }}
	}
}

package metrics

import "sync"

func WithDriver(d func() Driver, bulkSize int, concurrency int) Option {
	driver := d()
	return func(r *Reporter) {
		r.buff = &sync.Pool{New: func() interface{} { return newRequestBuffer(bulkSize, driver) }}
	}
}

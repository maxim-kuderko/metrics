package metrics

import "time"

func WithDriver(d Driver) Option {
	return func(r *Reporter) {
		r.driver = d
	}
}

func WithBuffer(size int) Option {
	return func(r *Reporter) {
		r.buffSize = size
	}
}

func WithFlushTicker(duration time.Duration) Option {
	return func(r *Reporter) {
		r.bufferFlushTicker = time.NewTicker(duration)
		r.flushTickerDuration = duration
	}
}

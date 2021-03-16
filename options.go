package metrics

func WithDriver(d func() Driver, bulkSize int, concurrency int) Option {
	return func(r *Reporter) {
		requestBuffers := make([]*requestBuffer, 0, concurrency)
		for i := 0; i < concurrency; i++ {
			dd := d()
			requestBuffers = append(requestBuffers, newRequestBuffer(bulkSize, dd))
		}
		r.buff = requestBuffers
	}
}

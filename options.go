package metrics

func WithDriver(d func() Driver, concurrency int) Option {
	return func(r *Reporter) {
		drivers := make([]Driver, 0, concurrency)
		for i := 0; i < concurrency; i++ {
			drivers = append(drivers, d())
		}
		r.driver = drivers
	}
}

package metrics

import "github.com/maxim-kuderko/metrics/drivers"

type Driver interface {
	Send(metrics drivers.Metrics)
}

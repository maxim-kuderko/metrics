package metrics

import (
	"github.com/maxim-kuderko/metrics/entities"
)

type Driver interface {
	Send(metrics entities.Metrics)
}

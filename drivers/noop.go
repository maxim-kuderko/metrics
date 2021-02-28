package drivers

import "github.com/maxim-kuderko/metrics/entities"

type Noop struct {
}

func (s Noop) Send(metrics entities.Metrics) {
}

func NewNoop() *Noop {
	return &Noop{}
}

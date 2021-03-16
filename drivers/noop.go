package drivers

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
)

type Noop struct {
}

func (s Noop) Send(metrics *proto.Metric) {
}
func (s Noop) Close() {
}

func NewNoop() *Noop {
	return &Noop{}
}

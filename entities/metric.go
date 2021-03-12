package entities

import "github.com/maxim-kuderko/metrics-collector/proto"

type Metrics map[uint64]*proto.Metric

func (m Metrics) Reset() {
	for k := range m {
		delete(m, k)
	}
}

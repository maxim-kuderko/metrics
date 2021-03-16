package metrics

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
)

type Driver interface {
	Send(metrics *proto.MetricsRequest)
	Close()
}

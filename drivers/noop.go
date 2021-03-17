package drivers

import (
	marshaler "github.com/golang/protobuf/proto"
	"github.com/maxim-kuderko/metrics-collector/proto"
)

type Noop struct {
}

func (s Noop) Send(metrics *proto.MetricsRequest) {
	buff := marshlerPool.Get().(*marshaler.Buffer)
	defer func() {
		buff.Reset()
		marshlerPool.Put(buff)
	}()
	buff.Marshal(metrics)
}
func (s Noop) Close() {
}

func NewNoop() *Noop {
	return &Noop{}
}

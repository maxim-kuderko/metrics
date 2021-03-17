package drivers

import (
	marshaler "github.com/golang/protobuf/proto"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"io"
	"net"
	"sync"
)

type UDP struct {
	w io.Writer
}

const UDPBufferSize = 8 << 13

var marshlerPool = &sync.Pool{New: func() interface{} {
	return marshaler.NewBuffer(nil)
}}

func (s *UDP) Send(metrics *proto.MetricsRequest) {
	buff := marshlerPool.Get().(*marshaler.Buffer)
	defer func() {
		buff.Reset()
		marshlerPool.Put(buff)
	}()
	buff.Marshal(metrics)
	s.w.Write(buff.Bytes())
	counter.Send(metrics)
}

func (s *UDP) Close() {
}

func NewUDP(addr string) *UDP {
	NewCounter()
	c, err := net.Dial(`udp`, addr)
	if err != nil {
		panic(err)
	}
	return &UDP{
		w: c,
	}
}

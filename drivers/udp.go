package drivers

import (
	marshaler "github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"net"
	"sync"
)

type UDP struct {
	c net.Conn
	w *snappy.Writer
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
	s.w.Flush()
	s.w.Reset(s.c)
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
		c: c,
		w: snappy.NewBufferedWriter(c),
	}
}

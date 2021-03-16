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

	mu sync.Mutex
}

const UDPBufferSize = 8 << 10

func (s *UDP) Send(metrics *proto.MetricsRequest) {
	b, _ := marshaler.Marshal(metrics)
	s.w.Write(b)
}

func (s *UDP) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

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

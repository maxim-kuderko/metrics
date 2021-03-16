package drivers

import (
	"bytes"
	"encoding/binary"
	marshaler "github.com/golang/protobuf/proto"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"io"
	"net"
	"sync"
)

type UDP struct {
	w    io.Writer
	buff *bytes.Buffer

	mu sync.Mutex
}

const UDPBufferSize = 8 << 9

func (s *UDP) Send(metrics *proto.Metric) {
	b, _ := marshaler.Marshal(metrics)
	si := make([]byte, 4)
	binary.BigEndian.PutUint32(si, uint32(len(b)))
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.buff.Len()+len(b)+len(si) >= UDPBufferSize {
		s.buff.WriteTo(s.w)
		s.buff.Reset()
	}
	s.buff.Write(si)
	s.buff.Write(b)
	counter.Send(metrics)
}

func NewUDP(addr string) *UDP {
	NewCounter()
	c, err := net.Dial(`udp`, addr)
	if err != nil {
		panic(err)
	}
	return &UDP{
		w:    c,
		buff: bytes.NewBuffer(nil),
	}
}

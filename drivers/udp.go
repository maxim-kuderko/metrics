package drivers

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"net"
)

type UDP struct {
	conn net.Conn
	enc  *jsoniter.Encoder
}

func (s *UDP) Send(metrics *proto.Metric) {
	s.enc.Encode(metrics)
}

func NewUDP(addr string) *UDP {
	c, err := net.Dial(`udp`, addr)
	if err != nil {
		panic(err)
	}
	return &UDP{
		conn: c,
		enc:  jsoniter.ConfigFastest.NewEncoder(c),
	}
}

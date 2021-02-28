package drivers

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/entities"
	"github.com/valyala/bytebufferpool"
	"net"
	"sync"
)

type TCP struct {
	conn net.Conn
	mu   sync.Mutex
}

func (s *TCP) Send(metrics entities.Metrics) {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	enc := jsoniter.ConfigFastest.NewEncoder(b)
	for _, m := range metrics {
		enc.Encode(m)
		if b.Len() > udpBufferSize/4 {
			s.flush(b)
			b.Reset()
		}
	}
	if b.Len() > 0 {
		s.flush(b)
	}
}
func (s *TCP) flush(buffer *bytebufferpool.ByteBuffer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	buffer.WriteTo(s.conn)
}

func NewTCP(addr string) (*TCP, error) {
	c, err := net.Dial(`tcp`, addr)
	return &TCP{
		conn: c,
	}, err
}

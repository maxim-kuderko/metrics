package drivers

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/entities"
	"github.com/valyala/bytebufferpool"
	"net"
	"sync"
)

type UDP struct {
	conn net.Conn
	mu   sync.Mutex
}

const udpBufferSize = 6500

func (s *UDP) Send(metrics entities.Metrics) {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	enc := jsoniter.ConfigFastest.NewEncoder(b)
	for _, m := range metrics {
		enc.Encode(m)
		if b.Len() > udpBufferSize {
			s.flush(b)
			b.Reset()
		}
	}
	if b.Len() > 0 {
		s.flush(b)
	}
}
func (s *UDP) flush(buffer *bytebufferpool.ByteBuffer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	buffer.WriteTo(s.conn)
}

func NewUDP(addr string) (*UDP, error) {
	c, err := net.Dial(`udp`, addr)
	return &UDP{
		conn: c,
	}, err
}

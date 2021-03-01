package drivers

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/entities"
	"net/http"
	"sync"
	"time"
)

type HTTP struct {
	conn http.Client
	addr string
	mu   sync.Mutex
}

var bytesBuffer = sync.Pool{New: func() interface{} { return bytes.NewBuffer(nil) }}

func (s *HTTP) Send(metrics entities.Metrics) {
	b := bytesBuffer.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		bytesBuffer.Put(b)
	}()
	enc := jsoniter.ConfigFastest.NewEncoder(b)
	for _, m := range metrics {
		enc.Encode(m)
	}
	s.flush(b)
}
func (s *HTTP) flush(buffer *bytes.Buffer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resp, err := s.conn.Post(s.addr, ``, buffer)
	if err != nil {
		fmt.Println(`error sending http `, err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		fmt.Println(`error seding http status code `, resp.StatusCode)
	}
}

func NewHTTP(addr string, timeout time.Duration) *HTTP {
	return &HTTP{
		conn: http.Client{
			Timeout: timeout,
		},
		addr: addr,
	}
}

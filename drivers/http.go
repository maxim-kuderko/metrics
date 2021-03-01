package drivers

import (
	"bytes"
	"fmt"
	"github.com/golang/snappy"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/entities"
	"io"
	"io/ioutil"
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
	r, w := io.Pipe()
	sn := snappy.NewBufferedWriter(w)
	enc := jsoniter.ConfigFastest.NewEncoder(sn)
	go s.flush(r)
	defer func() {
		sn.Close()
		w.Close()
	}()
	for _, m := range metrics {
		enc.Encode(m)
	}

}
func (s *HTTP) flush(buffer io.Reader) {
	resp, err := http.Post(s.addr, ``, buffer)
	if err != nil {
		io.Copy(ioutil.Discard, buffer)
		fmt.Println(`error sending http `, err)
		return
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

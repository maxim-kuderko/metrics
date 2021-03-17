package drivers

import (
	"compress/gzip"
	"fmt"
	marshaler "github.com/golang/protobuf/proto"
	"github.com/klauspost/compress/snappy"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/valyala/bytebufferpool"
	"net"
	"sync"
)

type UDP struct {
	c net.Conn
	w *gzip.Writer
}

const UDPBufferSize = 8 << 10

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
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	w := snappy.NewWriter(b)
	w.Write(buff.Bytes())
	w.Close()
	if b.Len() > UDPBufferSize {
		fmt.Println(`bbbb`)
	}
	b.WriteTo(s.c)
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
	}
}

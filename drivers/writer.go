package drivers

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"io"
)

type Writer struct {
	w io.Writer
}

func (s Writer) Send(metrics *proto.Metric) {
	jsoniter.ConfigFastest.NewEncoder(s.w).Encode(metrics)
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

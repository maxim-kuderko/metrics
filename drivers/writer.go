package drivers

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/entities"
	"io"
)

type Writer struct {
	w io.Writer
}

func (s Writer) Send(metrics entities.Metrics) {
	for _, m := range metrics {
		jsoniter.ConfigFastest.NewEncoder(s.w).Encode(m)
	}
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

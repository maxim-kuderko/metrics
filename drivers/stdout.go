package drivers

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics/entities"
	"os"
)

type Stdout struct {
}

func (s Stdout) Send(metrics entities.Metrics) {
	for _, m := range metrics {
		jsoniter.ConfigFastest.NewEncoder(os.Stdout).Encode(m)
	}
}

func NewStdout() *Stdout {
	return &Stdout{}
}

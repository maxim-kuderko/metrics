package drivers

import (
	"fmt"
	"github.com/maxim-kuderko/metrics/entities"
)

type Stdout struct {
}

func (s Stdout) Send(metrics entities.Metrics) {
	for _, m := range metrics {
		fmt.Println(m)
	}
}

func NewStdout() *Stdout {
	return &Stdout{}
}

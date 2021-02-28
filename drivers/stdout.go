package drivers

import (
	"fmt"
)

type Stdout struct {
}

func (s Stdout) Send(metrics Metrics) {
	for _, m := range metrics {
		fmt.Println(m)
	}
}

func NewStdout() *Stdout {
	return &Stdout{}
}

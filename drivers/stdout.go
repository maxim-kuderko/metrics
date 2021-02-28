package drivers

import (
	"fmt"
	"github.com/maxim-kuderko/metrics/entities"
)

type Stdout struct {
}

func (s Stdout) Send(metrics entities.Metrics) {
	fmt.Println(metrics)
}

func NewStdout() *Stdout {
	return &Stdout{}
}

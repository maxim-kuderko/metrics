package drivers

type Noop struct {
}

func (s Noop) Send(metrics Metrics) {
}

func NewNoop() *Noop {
	return &Noop{}
}

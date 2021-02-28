package drivers

type Metric struct {
	Name  string
	Value float64
	Tags  []string

	hash uint32
}

type Metrics []Metric

func (m *Metric) Hash() uint32 {
	return 0
}

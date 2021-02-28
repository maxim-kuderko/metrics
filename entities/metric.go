package entities

type Metrics map[uint64]*AggregatedMetric

func (m Metrics) Reset() {
	for k := range m {
		delete(m, k)
	}
}

type AggregatedMetric struct {
	Name   string
	Values Values
	Tags   []string
}
type Values struct {
	Count int64
	Sum   float64
}

func (am *AggregatedMetric) Add(value float64) {
	am.Values.Count++
	am.Values.Sum += value
}

package entities

type Metrics map[uint64]*AggregatedMetric

func (m Metrics) Reset() {
	for k := range m {
		delete(m, k)
	}
}

type AggregatedMetric struct {
	Name   string   `json:"name"`
	Values Values   `json:"values"`
	Tags   []string `json:"tags"`
	Hash   uint64   `json:"hash"`
	Time   int64    `json:"ts"`
}
type Values struct {
	Count int64   `json:"count"`
	Sum   float64 `json:"sum"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

func (am *AggregatedMetric) Add(value float64) {
	am.Values.Count++
	am.Values.Sum += value
	if value < am.Values.Min {
		am.Values.Min = value
	}
	if value > am.Values.Max {
		am.Values.Max = value
	}
}

func (am *AggregatedMetric) Merge(new *AggregatedMetric) {
	am.Values.Count += new.Values.Count
	am.Values.Sum += new.Values.Sum
	if new.Values.Min < am.Values.Min {
		am.Values.Min = new.Values.Min
	}
	if new.Values.Max > am.Values.Max {
		am.Values.Max = new.Values.Max
	}
}

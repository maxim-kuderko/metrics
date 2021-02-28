package drivers

import (
	"github.com/maxim-kuderko/metrics/entities"
	"sync"
)

type TestStub struct {
	m  []*entities.AggregatedMetric
	mu sync.Mutex
}

func (s *TestStub) Send(metrics entities.Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range metrics {
		s.m = append(s.m, v)
	}
}
func (s *TestStub) Metrics() []*entities.AggregatedMetric {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m
}

func NewTestStub() *TestStub {
	return &TestStub{m: make([]*entities.AggregatedMetric, 0)}
}

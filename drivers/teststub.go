package drivers

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/maxim-kuderko/metrics/entities"
	"sync"
)

type TestStub struct {
	m  []*proto.Metric
	mu sync.Mutex
}

func (s *TestStub) Send(metrics entities.Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range metrics {
		s.m = append(s.m, v)
	}
}
func (s *TestStub) Metrics() []*proto.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m
}

func NewTestStub() *TestStub {
	return &TestStub{m: make([]*proto.Metric, 0)}
}

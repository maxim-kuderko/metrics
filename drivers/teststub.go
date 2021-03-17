package drivers

import (
	"github.com/maxim-kuderko/metrics-collector/proto"
	"sync"
)

type TestStub struct {
	m  []*proto.Metric
	mu sync.Mutex
}

func (s *TestStub) Send(metric *proto.MetricsRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = append(s.m, metric.Metrics...)

}
func (s *TestStub) Metrics() []*proto.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m
}

func NewTestStub() *TestStub {
	return &TestStub{m: make([]*proto.Metric, 0)}
}

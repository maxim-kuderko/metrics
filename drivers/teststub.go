package drivers

import (
	"github.com/maxim-kuderko/metrics/entities"
	"sync"
)

type TestStub struct {
	m  entities.Metrics
	mu sync.Mutex
}

func (s *TestStub) Send(metrics entities.Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = append(s.m, metrics...)
}
func (s *TestStub) Metrics() entities.Metrics {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m
}

func NewTestStub() *TestStub {
	return &TestStub{m: make(entities.Metrics, 0)}
}

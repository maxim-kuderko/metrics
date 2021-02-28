package drivers

import (
	"sync"
)

type TestStub struct {
	m  Metrics
	mu sync.Mutex
}

func (s *TestStub) Send(metrics Metrics) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = append(s.m, metrics...)
}
func (s *TestStub) Metrics() Metrics {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m
}

func NewTestStub() *TestStub {
	return &TestStub{m: make(Metrics, 0)}
}

package drivers

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"google.golang.org/grpc"
	"sync"
)

type Grpc struct {
	c    proto.MetricsCollectorGrpcClient
	buff *proto.MetricsRequest
	mu   sync.Mutex
}

func (s *Grpc) Send(metrics *proto.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buff.Metric = append(s.buff.Metric, metrics)
	if len(s.buff.Metric) == cap(s.buff.Metric) {
		if _, err := s.c.Send(context.Background(), s.buff); err != nil {
			fmt.Println(err)
		}
		s.buff.Metric = s.buff.Metric[:0]
	}

}

func NewGrpc(ctx context.Context, url string, options ...grpc.DialOption) *Grpc {
	conn, err := grpc.DialContext(ctx, url, options...)
	if err != nil {
		panic(err)
	}

	c := proto.NewMetricsCollectorGrpcClient(conn)
	return &Grpc{c: c, buff: &proto.MetricsRequest{Metric: make([]*proto.Metric, 0, 10000)}}
}

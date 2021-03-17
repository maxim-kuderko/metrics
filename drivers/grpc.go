package drivers

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"google.golang.org/grpc"
	"io"
)

type Grpc struct {
	c proto.MetricsCollectorGrpcClient
}

func (s *Grpc) Send(metrics *proto.MetricsRequest) {
	counter.Send(metrics)
	if _, err := s.c.Bulk(context.Background(), metrics); err != nil && err != io.EOF {
		fmt.Println(err)
	}
}

func NewGrpc(ctx context.Context, url string, options ...grpc.DialOption) *Grpc {
	NewCounter()
	conn, err := grpc.DialContext(ctx, url, options...)
	if err != nil {
		panic(err)
	}

	c := proto.NewMetricsCollectorGrpcClient(conn)
	return &Grpc{c: c}
}

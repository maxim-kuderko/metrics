package drivers

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"google.golang.org/grpc"
)

type Grpc struct {
	c proto.MetricsCollectorGrpcClient
}

func (s Grpc) Send(metrics *proto.MetricsRequest) {
	if _, err := s.c.Send(context.Background(), metrics); err != nil {
		fmt.Println(err)
	}
}

func NewGrpc(ctx context.Context, url string, options ...grpc.DialOption) *Grpc {
	conn, err := grpc.DialContext(ctx, url, options...)
	if err != nil {
		panic(err)
	}

	c := proto.NewMetricsCollectorGrpcClient(conn)
	return &Grpc{c: c}
}

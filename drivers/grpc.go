package drivers

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"google.golang.org/grpc"
	"io"
)

type Grpc struct {
	c proto.MetricsCollectorGrpc_SendClient
}

func (s *Grpc) Send(metrics *proto.Metric) {
	if err := s.c.Send(metrics); err != nil && err != io.EOF {
		fmt.Println(err)
	}
}
func (s Grpc) Close() {
	s.c.CloseSend()
}

func NewGrpc(ctx context.Context, url string, options ...grpc.DialOption) *Grpc {
	conn, err := grpc.DialContext(ctx, url, options...)
	if err != nil {
		panic(err)
	}

	c, err := proto.NewMetricsCollectorGrpcClient(conn).Send(context.Background())
	if err != nil {
		panic(err)
	}
	return &Grpc{c: c}
}

package drivers

import (
	"context"
	"fmt"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/maxim-kuderko/metrics/entities"
	"google.golang.org/grpc"
	"io"
)

type Grpc struct {
	c proto.MetricsCollectorGrpc_SendClient
}

func (s Grpc) Send(metrics entities.Metrics) {
	for _, m := range metrics {
		if err := s.c.Send(m); err != nil && err != io.EOF {
			fmt.Println(err)
		}
	}
}

func NewGrpc(url string, options ...grpc.DialOption) *Grpc {
	conn, err := grpc.Dial(url, options...)
	if err != nil {
		panic(err)
	}

	c, err := proto.NewMetricsCollectorGrpcClient(conn).Send(context.Background())
	if err != nil {
		panic(err)
	}
	return &Grpc{c: c}
}

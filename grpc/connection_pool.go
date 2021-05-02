package goload_grpc

import (
	"log"
	"math/rand"

	"google.golang.org/grpc"
)

type GRPCConnectionPool struct {
	connections []*grpc.ClientConn
}

func NewGRPCConnectionPool(
	count int,
	target string,
	opts ...grpc.DialOption,
) *GRPCConnectionPool {
	connections := make([]*grpc.ClientConn, count)

	for i := 0; i < count; i++ {
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			log.Fatal(err)
		}

		connections[i] = conn
	}

	return &GRPCConnectionPool{
		connections: connections,
	}
}

func (pool *GRPCConnectionPool) GetConnection() *grpc.ClientConn {
	return pool.connections[rand.Intn(len(pool.connections))]
}

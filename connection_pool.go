package goload

import (
	"log"
	"math/rand"
	"net/http"

	"google.golang.org/grpc"
)

type HTTPClientPool struct {
	clients []*http.Client
}

func NewHTTPConnectionPool(count int) *HTTPClientPool {
	clients := make([]*http.Client, count)

	for i := 0; i < count; i++ {
		clients[i] = &http.Client{}
	}

	return &HTTPClientPool{
		clients: clients,
	}
}

func (pool *HTTPClientPool) GetClient() *http.Client {
	return pool.clients[rand.Intn(len(pool.clients))]
}

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

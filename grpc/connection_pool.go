package goload_grpc

import (
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

type ConnectionPool struct {
	conns []*grpc.ClientConn
	rand  *rand.Rand
}

func NewConnectionPool(
	count int,
	target string,
	opts ...grpc.DialOption,
) *ConnectionPool {
	conns := make([]*grpc.ClientConn, count)

	for i := 0; i < count; i++ {
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			log.Fatalf("Unable to dial GRPC connection: %v", err)
		}

		conns[i] = conn
	}

	return &ConnectionPool{
		conns: conns,
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (pool *ConnectionPool) Connection() *grpc.ClientConn {
	return pool.conns[pool.rand.Intn(len(pool.conns))]
}

func (pool *ConnectionPool) Close() {
	for _, conn := range pool.conns {
		conn.Close()
	}
}

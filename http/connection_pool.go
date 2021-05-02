package goload_http

import (
	"math/rand"
	"net/http"
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

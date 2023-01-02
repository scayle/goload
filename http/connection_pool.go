package goload_http

import (
	"math/rand"
	"net/http"
	"time"
)

type ClientPool struct {
	clients []*http.Client
	rand    *rand.Rand
}

type HTTPTransportOption = func(request *http.Request)

func WithUserAgent(userAgent string) HTTPTransportOption {
	return func(request *http.Request) {
		if request.Header.Get("user-agent") == "" {
			request.Header.Set("user-agent", userAgent)
		}
	}
}

type transport struct {
	options []HTTPTransportOption

	innerTransport http.RoundTripper
}

func (transport *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	for _, option := range transport.options {
		option(request)
	}

	return transport.innerTransport.RoundTrip(request)
}

func NewClientPool(
	count int,
	options ...HTTPTransportOption,
) *ClientPool {
	clients := make([]*http.Client, count)

	for i := 0; i < count; i++ {
		clients[i] = &http.Client{
			Transport: &transport{
				options:        options,
				innerTransport: http.DefaultTransport,
			},
		}
	}

	return &ClientPool{
		clients: clients,
		rand:    rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (pool *ClientPool) Client() *http.Client {
	return pool.clients[pool.rand.Intn(len(pool.clients))]
}

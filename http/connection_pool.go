package goload_http

import (
	"math/rand"
	"net/http"
)

type HTTPClientPool struct {
	clients []*http.Client
}

type HTTPTransportOption = func(request *http.Request)

func WithUserAgentOption(userAgent string) HTTPTransportOption {
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

func NewHTTPConnectionPool(
	count int,
	options ...HTTPTransportOption,
) *HTTPClientPool {
	clients := make([]*http.Client, count)

	for i := 0; i < count; i++ {
		clients[i] = &http.Client{
			Transport: &transport{
				options:        options,
				innerTransport: http.DefaultTransport,
			},
		}
	}

	return &HTTPClientPool{
		clients: clients,
	}
}

func (pool *HTTPClientPool) GetClient() *http.Client {
	return pool.clients[rand.Intn(len(pool.clients))]
}

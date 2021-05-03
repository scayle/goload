package goload_http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/HenriBeck/goload"
)

type HTTPMethod string

const (
	HTTPMethodGet  HTTPMethod = "GET"
	HTTPMethodPost HTTPMethod = "POST"
)

type HTTPEndpoint struct {
	Options *HTTPEndpointOptions
}

type HTTPEndpointOptions struct {
	URI *url.URL

	RequestsPerMinute int32

	Method HTTPMethod

	GetClient func() *http.Client
}

type HTTPEndpointConfig func(options *HTTPEndpointOptions)

func WithClientPool(pool *HTTPClientPool) HTTPEndpointConfig {
	return func(options *HTTPEndpointOptions) {
		options.GetClient = pool.GetClient
	}
}

func WithHTTPClient(client *http.Client) HTTPEndpointConfig {
	return func(options *HTTPEndpointOptions) {
		options.GetClient = func() *http.Client {
			return client
		}
	}
}

func WithHTTPMethod(method HTTPMethod) HTTPEndpointConfig {
	return func(options *HTTPEndpointOptions) {
		options.Method = method
	}
}

func WithRequestsPerMinute(rpm int32) HTTPEndpointConfig {
	return func(options *HTTPEndpointOptions) {
		options.RequestsPerMinute = rpm
	}
}

func WithURL(uri url.URL) HTTPEndpointConfig {
	return func(options *HTTPEndpointOptions) {
		options.URI = &uri
	}
}

func WithURLString(uri string) HTTPEndpointConfig {
	parsedUri, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	return func(options *HTTPEndpointOptions) {
		options.URI = parsedUri
	}
}

func NewHTTPEndpoint(
	configs ...HTTPEndpointConfig,
) goload.Endpoint {
	options := &HTTPEndpointOptions{
		Method: HTTPMethodGet,
		GetClient: func() *http.Client {
			return http.DefaultClient
		},
	}
	for _, config := range configs {
		config(options)
	}

	return &HTTPEndpoint{
		Options: options,
	}
}

func (endpoint *HTTPEndpoint) GetRequestsPerMinute() int32 {
	return endpoint.Options.RequestsPerMinute
}

func (endpoint *HTTPEndpoint) Execute(ctx context.Context) error {
	client := endpoint.Options.GetClient()

	var req *http.Request
	// TODO: more HTTP methods
	switch endpoint.Options.Method {
	case HTTPMethodGet:
		_req, err := http.NewRequestWithContext(
			ctx,
			"GET",
			endpoint.Options.URI.String(),
			nil,
		)
		if err != nil {
			return err
		}

		req = _req
	case HTTPMethodPost:
		_req, err := http.NewRequestWithContext(
			ctx,
			"POST",
			endpoint.Options.URI.String(),
			nil, // TODO: Body
		)
		if err != nil {
			return err
		}

		req = _req
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("received non 200 status code from request")
	}

	return nil
}

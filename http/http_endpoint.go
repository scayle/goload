package goload_http

import (
	"context"
	"fmt"
	"io"
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

	Method string

	Body io.Reader

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

func WithHTTPMethod(method string) HTTPEndpointConfig {
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
		Method: http.MethodGet,
		Body:   http.NoBody,
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

func (endpoint *HTTPEndpoint) Name() string {
	return fmt.Sprintf(
		"%s %s",
		endpoint.Options.Method,
		endpoint.Options.URI.String(),
	)
}

func (endpoint *HTTPEndpoint) Execute(ctx context.Context) error {
	client := endpoint.Options.GetClient()

	req, err := http.NewRequestWithContext(
		ctx,
		endpoint.Options.Method,
		endpoint.Options.URI.String(),
		endpoint.Options.Body,
	)
	if err != nil {
		return err
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

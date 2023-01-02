package goload_http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/HenriBeck/goload"
)

type Endpoint struct {
	options *EndpointOptions
}

type EndpointOptions struct {
	url *url.URL

	requestsPerMinute int32

	method string

	body io.Reader

	getClient func() *http.Client
}

type HTTPEndpointConfig func(options *EndpointOptions)

func WithClientPool(pool *ClientPool) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.getClient = pool.Client
	}
}

func WithHTTPClient(client *http.Client) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.getClient = func() *http.Client {
			return client
		}
	}
}

func WithHTTPMethod(method string) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.method = method
	}
}

func WithRequestsPerMinute(rpm int32) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.requestsPerMinute = rpm
	}
}

func WithURL(uri url.URL) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.url = &uri
	}
}

func WithURLString(uri string) HTTPEndpointConfig {
	parsedUri, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	return func(options *EndpointOptions) {
		options.url = parsedUri
	}
}

func NewEndpoint(configs ...HTTPEndpointConfig) goload.Endpoint {
	options := &EndpointOptions{
		method: http.MethodGet,
		body:   http.NoBody,
		getClient: func() *http.Client {
			return http.DefaultClient
		},
	}
	for _, config := range configs {
		config(options)
	}

	return &Endpoint{
		options: options,
	}
}

func (endpoint *Endpoint) GetRequestsPerMinute() int32 {
	return endpoint.options.requestsPerMinute
}

func (endpoint *Endpoint) Name() string {
	return fmt.Sprintf(
		"%s %s",
		endpoint.options.method,
		endpoint.options.url.String(),
	)
}

func (endpoint *Endpoint) Execute(ctx context.Context) error {
	client := endpoint.options.getClient()

	req, err := http.NewRequestWithContext(
		ctx,
		endpoint.options.method,
		endpoint.options.url.String(),
		endpoint.options.body,
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

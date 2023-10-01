package goload

import (
	"context"
	"time"
)

type EndpointOptions struct {
	RequestsPerMinute *int32
	Timeout           time.Duration
}

func (options *EndpointOptions) Apply(opts ...EndpointOption) {
	for _, opt := range opts {
		opt(options)
	}
}

type EndpointOption func(options *EndpointOptions)

func WithRequestsPerMinute(requestsPerMinute int32) EndpointOption {
	return func(options *EndpointOptions) {
		options.RequestsPerMinute = &requestsPerMinute
	}
}

func WithTimeout(timeout time.Duration) EndpointOption {
	return func(options *EndpointOptions) {
		options.Timeout = timeout
	}
}

type Endpoint interface {
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	Execute(ctx context.Context) error

	Name() string

	Options() *EndpointOptions
}

type endpoint struct {
	name    string
	handler func(ctx context.Context) error
	options *EndpointOptions
}

func (e *endpoint) Execute(ctx context.Context) error {
	return e.handler(ctx)
}

func (e *endpoint) Name() string {
	return e.name
}

func (e *endpoint) Options() *EndpointOptions {
	return e.options
}

func NewEndpoint(name string, handler func(ctx context.Context) error, opts ...EndpointOption) Endpoint {
	options := &EndpointOptions{
		Timeout:           0,
		RequestsPerMinute: nil,
	}
	options.Apply(opts...)

	return &endpoint{
		name:    name,
		handler: handler,
		options: options,
	}
}

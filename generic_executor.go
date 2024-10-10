package goload

import (
	"context"
	"time"
)

type genericExecutor struct {
	name    string
	handler func(ctx context.Context) error
	options *ExecutorOptions
}

type ExecutorOption func(options *ExecutorOptions)

func WithWeight(weight int) ExecutorOption {
	return func(options *ExecutorOptions) {
		options.Weight = weight
	}
}

func WithTimeout(timeout time.Duration) ExecutorOption {
	return func(options *ExecutorOptions) {
		options.Timeout = timeout
	}
}

func NewGenericExecutor(name string, handler func(ctx context.Context) error, opts ...ExecutorOption) Executor {
	options := &ExecutorOptions{
		Weight:  1,
		Timeout: 0,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &genericExecutor{
		name:    name,
		handler: handler,
		options: options,
	}
}

func (d *genericExecutor) Execute(ctx context.Context) ExecutionResponse {
	err := d.handler(ctx)
	return ExecutionResponse{
		Identifier:     d.Name(),
		Err:            err,
		AdditionalData: nil,
	}
}

func (d *genericExecutor) Name() string {
	return d.name
}

func (d *genericExecutor) Options() *ExecutorOptions {
	return d.options
}

package goload

import (
	"context"
	"time"
)

type ExecutionResponse struct {
	Identifier     string
	Err            error
	AdditionalData any
}

type ExecutorOptions struct {
	Weight  int
	Timeout time.Duration
}

type Executor interface {
	Execute(ctx context.Context) ExecutionResponse
	Name() string
	// TODO: remove pointer and maybe split into two functions
	Options() *ExecutorOptions
}

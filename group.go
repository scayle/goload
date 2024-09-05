package goload

import (
	"context"
	"github.com/mroth/weightedrand/v2"
	"time"
)

type executorGroup struct {
	name    string
	chooser *weightedrand.Chooser[Executor, int]
	weight  int
	timeout time.Duration
}

func (e *executorGroup) Execute(ctx context.Context) ExecutionResponse {
	return e.chooser.Pick().Execute(ctx)
}

func (e *executorGroup) Name() string {
	return e.name
}

func (e *executorGroup) Options() *ExecutorOptions {
	return &ExecutorOptions{
		Weight:  e.weight,
		Timeout: e.timeout,
	}
}

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

// TODO: add options
func NewGroup(name string, weight int, executors []Executor) Executor {
	if len(executors) == 0 {
		panic("group can't be empty")
	}
	choises := make([]weightedrand.Choice[Executor, int], 0, len(executors))
	for _, exec := range executors {
		choises = append(choises, weightedrand.NewChoice(exec, exec.Options().Weight))
	}
	chooser, err := weightedrand.NewChooser(
		choises...,
	)
	if err != nil {
		panic(err)
	}

	return &executorGroup{
		name:    name,
		chooser: chooser,
		weight:  weight,
		timeout: 0,
	}
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

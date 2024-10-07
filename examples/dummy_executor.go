package main

import (
	"context"
	"fmt"
	"github.com/scayle/goload"
)

type DummyExecutor struct {
	text string
	name string
}

func (d DummyExecutor) Execute(ctx context.Context) goload.ExecutionResponse {
	fmt.Println(d.text)
	return goload.ExecutionResponse{
		Identifier:     d.Name(),
		Err:            nil,
		AdditionalData: nil,
	}
}

func (d DummyExecutor) Name() string {
	return d.name
}

func (d DummyExecutor) Options() *goload.ExecutorOptions {
	return &goload.ExecutorOptions{
		Weight:  1,
		Timeout: 0,
	}
}

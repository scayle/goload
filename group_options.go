package goload

import (
	"github.com/mroth/weightedrand/v2"
	"github.com/rs/zerolog/log"
	"time"
)

type GroupOptions struct {
	name      string
	weight    int
	timeout   time.Duration
	executors []Executor
}

type GroupOption func(*GroupOptions)

func WithGroup(opts ...GroupOption) Executor {
	options := &GroupOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if len(options.executors) == 0 {
		log.Fatal().Msg("group can't be empty")
	}
	choises := make([]weightedrand.Choice[Executor, int], 0, len(options.executors))
	for _, exec := range options.executors {
		choises = append(choises, weightedrand.NewChoice(exec, exec.Options().Weight))
	}
	chooser, err := weightedrand.NewChooser(
		choises...,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create chooser")
	}

	if options.weight == 0 {
		weightSum := 0
		for _, exec := range options.executors {
			weightSum += exec.Options().Weight
		}
		options.weight = weightSum
	}

	return &executorGroup{
		name:    options.name,
		chooser: chooser,
		weight:  options.weight,
		timeout: options.timeout,
	}
}

func WithGroupWeight(weight int) GroupOption {
	return func(options *GroupOptions) {
		options.weight = weight
	}
}

func WithGroupTimeout(timeout time.Duration) GroupOption {
	return func(options *GroupOptions) {
		options.timeout = timeout
	}
}

func WithGroupExecutors(executors ...Executor) GroupOption {
	return func(options *GroupOptions) {
		options.executors = append(options.executors, executors...)
	}
}

func WithGroupName(name string) GroupOption {
	return func(options *GroupOptions) {
		options.name = name
	}
}

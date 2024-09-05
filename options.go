package goload

import (
	"context"
	"github.com/HenriBeck/goload/pacer"
	"time"
)

func WithDuration(loadTestDuration time.Duration) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.duration = loadTestDuration
	}
}

func WithExecutors(executors ...Executor) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.executors = append(options.executors, executors...)
	}
}

func WithLinearRampUpPacer(startRate, endRate pacer.Rate, rampUpDuration time.Duration) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.pacer = pacer.NewLinearRampUpPacer(startRate, endRate, rampUpDuration)
	}
}

func WithConstantPacer(rate pacer.Rate) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.pacer = pacer.NewConstantPacer(rate)
	}
}

func WithReportInterval(reportInterval time.Duration) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.reportInterval = reportInterval
	}
}

func WithInitialWorkerCount(count int) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.initialWorkers = count
	}
}

func WithMaxWorkerCount(count int) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.maxWorkers = count
	}
}

func WithAdditionalResultHandler(handler resultHandler) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.resultHandlers = append(options.resultHandlers, handler)
	}
}

func WithContextModifier(fn func(ctx context.Context) context.Context) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.ctxModifier = fn
	}
}

func WithDefaultTimeout(timeout time.Duration) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.defaultTimeout = timeout
	}
}

func WithWeightOverrides(overrides map[string]int) LoadTestOption {
	return func(options *LoadTestOptions) {
		options.weightOverrides = overrides
	}
}

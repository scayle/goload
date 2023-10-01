package goload

import (
	"context"
	"time"
)

type LoadTestOptions struct {
	Endpoints []Endpoint
	// A config which allows to set the requests per minute for each registered endpoint at a central place
	//
	// This overrides any local endpoint setting where a default RPM value was set
	RequestPerMinutePerEndpoint map[string]int32

	DefaultEndpointTimeout time.Duration

	LoadTestDuration time.Duration

	RPMStrategy RPMStrategy

	ContextModifiers []func(ctx context.Context) context.Context
}

type LoadTestConfig func(options *LoadTestOptions)

func WithEndpoints(endpoints ...Endpoint) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.Endpoints = append(options.Endpoints, endpoints...)
	}
}

func WithRequestsPerMinuteForEndpoints(data map[string]int32) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.RequestPerMinutePerEndpoint = data
	}
}

func WithContextModifier(fn func(ctx context.Context) context.Context) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.ContextModifiers = append(options.ContextModifiers, fn)
	}
}

func WithDefaultEndpointTimeout(timeout time.Duration) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.DefaultEndpointTimeout = timeout
	}
}

func WithDuration(loadTestDuration time.Duration) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.LoadTestDuration = loadTestDuration
	}
}

func WithRPMStrategy(strategy RPMStrategy) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.RPMStrategy = strategy
	}
}

func WithStaticRPM(rpm int32) LoadTestConfig {
	return WithRPMStrategy(NewStaticRPMStrategy(rpm))
}

func WithRampUpRPM(steps []Step) LoadTestConfig {
	return WithRPMStrategy(NewRampUpRPMStrategy(steps))
}

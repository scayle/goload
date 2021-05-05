package goload

import (
	"time"
)

type LoadTestOptions struct {
	Endpoints []Endpoint

	DefaultEndpointTimeout time.Duration

	LoadTestDuration time.Duration

	RPMStrategy RPMStrategy
}

type LoadTestConfig func(options *LoadTestOptions)

func WithEndpoint(endpoint Endpoint) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.Endpoints = append(options.Endpoints, endpoint)
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

func WithStaticRPM(rpm int32) LoadTestConfig {
	strategy := NewStaticRPMStrategy(rpm)

	return func(options *LoadTestOptions) {
		options.RPMStrategy = strategy
	}
}

func WithRampUpRPM(steps []Step) LoadTestConfig {
	strategy := NewRampUpRPMStrategy(steps)

	return func(options *LoadTestOptions) {
		options.RPMStrategy = strategy
	}
}

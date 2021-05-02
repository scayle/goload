package goload

import "time"

type LoadTestOptions struct {
	Endpoints []Endpoint

	DefaultEndpointTimeout time.Duration

	LoadTestDuration time.Duration

	RequestsPerMinutePerNode uint
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

func WithRequestsPerMinutePerNode(rpm uint) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.RequestsPerMinutePerNode = rpm
	}
}

func (options *LoadTestOptions) GetRequestsPerMinute() uint {
	if options.RequestsPerMinutePerNode > 0 {
		return options.RequestsPerMinutePerNode
	}

	totalRequestsPerMinute := uint(0)
	for _, endpoint := range options.Endpoints {
		totalRequestsPerMinute += endpoint.GetRequestsPerMinute()
	}

	return totalRequestsPerMinute
}

package goload

import (
	"context"
	"time"
)

type Endpoint interface {
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	Execute(ctx context.Context) error

	GetRequestsPerMinute() int32

	Name() string
}

type EndpointWithTimeout interface {
	GetTimeout() time.Duration
}

func getTimeoutForEndpoint(endpoint Endpoint, options *LoadTestOptions) (time.Duration, bool) {
	endpointWithTimeout, ok := endpoint.(EndpointWithTimeout)
	if ok {
		return endpointWithTimeout.GetTimeout(), true
	}

	if options.DefaultEndpointTimeout.Nanoseconds() > 0 {
		return options.DefaultEndpointTimeout, true
	}

	return 0, false
}

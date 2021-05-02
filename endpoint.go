package goload

import (
	"context"
	"time"

	"github.com/mroth/weightedrand"
)

type Endpoint interface {
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	Execute(ctx context.Context) error

	GetRequestsPerMinute() uint
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

type EndpointWithTimeout interface {
	GetTimeout() time.Duration
}

type EndpointRandomizer struct {
	chooser *weightedrand.Chooser
}

func NewEndpointRandomizer(endpoints []Endpoint) *EndpointRandomizer {
	choices := make([]weightedrand.Choice, len(endpoints))
	for i, endpoint := range endpoints {
		choices[i] = weightedrand.NewChoice(
			endpoint,
			uint(endpoint.GetRequestsPerMinute()),
		)
	}

	chooser, _ := weightedrand.NewChooser(choices...)
	return &EndpointRandomizer{
		chooser: chooser,
	}
}

func (r *EndpointRandomizer) PickRandomEndpoint() Endpoint {
	return r.chooser.Pick().(Endpoint)
}

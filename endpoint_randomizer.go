package goload

import "github.com/mroth/weightedrand"

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

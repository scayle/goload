package goload

import (
	"fmt"
	"math/rand"
	"time"
)

type randomizedEndpoint struct {
	start, end int32
	endpoint   Endpoint
}

type EndpointRandomizer struct {
	total     int32
	endpoints []randomizedEndpoint
	rand      *rand.Rand
}

func NewEndpointRandomizer(endpoints []Endpoint, overrides map[string]int32) (*EndpointRandomizer, error) {
	randomizedEndpoints := make([]randomizedEndpoint, len(endpoints))
	var total int32

	for i, endpoint := range endpoints {
		var requestPerMinute int32

		if rpmOverride, ok := overrides[endpoint.Name()]; ok {
			requestPerMinute = rpmOverride
		} else if endpoint.Options().RequestsPerMinute != nil {
			requestPerMinute = *endpoint.Options().RequestsPerMinute
		} else {
			return nil, fmt.Errorf("Missing request per minute config for endpoint: %s", endpoint.Name())
		}

		randomizedEndpoints[i] = randomizedEndpoint{
			start:    total + 1,
			end:      total + requestPerMinute,
			endpoint: endpoint,
		}

		total += requestPerMinute
	}

	return &EndpointRandomizer{
		endpoints: randomizedEndpoints,
		total:     total,
		rand:      rand.New(rand.NewSource(time.Now().Unix())),
	}, nil
}

func (r *EndpointRandomizer) PickRandomEndpoint() Endpoint {
	pickedRange := r.rand.Int31n(r.total) + 1

	for _, endpoint := range r.endpoints {
		if endpoint.start <= pickedRange && pickedRange <= endpoint.end {
			return endpoint.endpoint
		}
	}

	return nil
}

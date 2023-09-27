package goload

import (
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

func NewEndpointRandomizer(endpoints []Endpoint) *EndpointRandomizer {
	randomizedEndpoints := make([]randomizedEndpoint, len(endpoints))
	var total int32

	for i, endpoint := range endpoints {
		randomizedEndpoints[i] = randomizedEndpoint{
			start:    total + 1,
			end:      total + endpoint.GetRequestsPerMinute(),
			endpoint: endpoint,
		}

		total += endpoint.GetRequestsPerMinute()
	}

	return &EndpointRandomizer{
		endpoints: randomizedEndpoints,
		total:     total,
		rand:      rand.New(rand.NewSource(time.Now().Unix())),
	}
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

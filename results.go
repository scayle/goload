package goload

import (
	"sort"
	"sync/atomic"
	"time"
)

type EndpointResult struct {
	Failed   bool
	Duration time.Duration
}

type EndpointResults struct {
	Name          string
	total         uint64
	failed        uint64
	totalDuration uint64
}

func (e *EndpointResults) GetTotalRequests() uint64 {
	return atomic.LoadUint64(&e.total)
}

func (e *EndpointResults) GetTotalFailedRequests() uint64 {
	return atomic.LoadUint64(&e.failed)
}

func (e *EndpointResults) GetAverageDuration() float64 {
	return float64(atomic.LoadUint64(&e.totalDuration)) / float64(atomic.LoadUint64(&e.total))
}

type LoadTestResults struct {
	endpoints map[string]*EndpointResults
}

func NewResults(endpoints []Endpoint) *LoadTestResults {
	results := &LoadTestResults{
		endpoints: map[string]*EndpointResults{},
	}

	for _, endpoint := range endpoints {
		results.endpoints[endpoint.Name()] = &EndpointResults{
			Name: endpoint.Name(),
		}
	}

	return results
}

func (results *LoadTestResults) Iter() []*EndpointResults {
	endpoints := []*EndpointResults{}
	for _, endpoint := range results.endpoints {
		endpoints = append(endpoints, endpoint)
	}

	sort.SliceStable(endpoints, func(i, j int) bool {
		return endpoints[i].Name < endpoints[j].Name
	})

	return endpoints
}

func (results *LoadTestResults) SaveEndpointResult(Endpoint Endpoint, result EndpointResult) {
	atomic.AddUint64(&results.endpoints[Endpoint.Name()].total, 1)
	if result.Failed {
		atomic.AddUint64(&results.endpoints[Endpoint.Name()].failed, 1)
	}

	atomic.AddUint64(
		&results.endpoints[Endpoint.Name()].totalDuration,
		uint64(result.Duration.Milliseconds()),
	)
}

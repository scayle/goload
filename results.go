package goload

import (
	"sort"
	"sync"
	"time"

	"github.com/spenczar/tdigest"
)

type EndpointResult struct {
	Failed   bool
	Duration time.Duration
}

type EndpointResults struct {
	Name    string
	lock    sync.Mutex
	results []EndpointResult
	td      *tdigest.TDigest
}

func (e *EndpointResults) GetTotalRequests() int {
	return len(e.results)
}

func (e *EndpointResults) GetTotalFailedRequests() int {
	c := 0
	for _, r := range e.results {
		if r.Failed {
			c++
		}
	}

	return c
}

func (e *EndpointResults) GetAverageDuration() float64 {
	totalDuration := float64(0)
	for _, r := range e.results {
		totalDuration += float64(r.Duration.Milliseconds())
	}

	return totalDuration / float64(len(e.results))
}

func (e *EndpointResults) GetPercentile(p float64) float64 {
	return e.td.Quantile(p)
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
			td:   tdigest.New(),
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
	e := results.endpoints[Endpoint.Name()]

	e.lock.Lock()
	e.results = append(e.results, result)
	e.td.Add(float64(result.Duration.Milliseconds()), 1)
	e.lock.Unlock()
}

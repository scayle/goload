package goload

import (
	"github.com/paulbellamy/ratecounter"
	"sync/atomic"
	"time"
)

type Result struct {
	Identifier     string
	Timestamp      time.Time
	Latency        time.Duration
	Err            error
	AdditionalData any
}

type resultAggregator struct {
	rateCounter *ratecounter.RateCounter
	total       atomic.Int64
	failures    atomic.Int64
}

func newResultAggregator() *resultAggregator {
	return &resultAggregator{
		rateCounter: ratecounter.NewRateCounter(10 * time.Second),
	}
}

func (ra *resultAggregator) resultAggregationHandler(_ *LoadTest, result *Result) {
	ra.rateCounter.Incr(1)
	ra.total.Add(1)
	if result.Err != nil {
		ra.failures.Add(1)
	}
}

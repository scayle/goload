package pacer

import (
	"math"
	"time"
)

func NewConstantPacer(rate Rate) Pacer {
	return constantPacer{
		rate: rate,
	}
}

type constantPacer struct {
	rate Rate
}

func (p constantPacer) Pace(elapsed time.Duration, hits uint64) time.Duration {
	expectedHits := p.expectedHits(elapsed)
	if hits == 0 || hits < uint64(expectedHits) {
		// Running behind, send next hit immediately.
		return 0
	}

	rate := p.Rate(elapsed)
	interval := math.Round(1e9 / rate)

	if n := uint64(interval); n != 0 && math.MaxInt64/n < hits {
		// We would overflow wait if we continued, so stop the execution.
		return 0
	}

	delta := float64(hits+1) - expectedHits
	wait := time.Duration(interval * delta)

	return wait
}

func (p constantPacer) Rate(elapsed time.Duration) float64 {
	return p.rate.hitsPerSec()
}

func (p constantPacer) expectedHits(t time.Duration) float64 {
	if t < 0 {
		return 0
	}

	return p.rate.hitsPerSec() * t.Seconds()
}

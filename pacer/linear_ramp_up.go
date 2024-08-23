package pacer

import (
	"math"
	"time"
)

// NewLinearRampUpPacer creates a linear increase from the start to Target rate within the given RampUpDuration
func NewLinearRampUpPacer(startRate Rate, targetRate Rate, rampUpDuration time.Duration) Pacer {
	diff := targetRate.hitsPerSec() - startRate.hitsPerSec()

	return linearRampUpPacer{
		StartRate:      startRate,
		TargetRate:     targetRate,
		RampUpDuration: rampUpDuration,
		a:              diff / rampUpDuration.Seconds(),
		b:              startRate.hitsPerSec(),
	}
}

type linearRampUpPacer struct {
	StartRate      Rate
	TargetRate     Rate
	RampUpDuration time.Duration

	a float64
	b float64
}

// Pace determines the length of time to sleep until the next hit is sent.
func (p linearRampUpPacer) Pace(elapsed time.Duration, hits uint64) time.Duration {
	switch {
	case p.StartRate.Per == 0 || p.StartRate.Freq == 0:
		return 0 // Zero value = infinite rate
	case p.StartRate.Per < 0 || p.StartRate.Freq < 0:
		panic("rate can't be negative")
	}

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

func (p linearRampUpPacer) Rate(elapsed time.Duration) float64 {
	if elapsed > p.RampUpDuration {
		return p.TargetRate.hitsPerSec()
	}

	x := elapsed.Seconds()

	return p.a*x + p.b
}

func (p linearRampUpPacer) expectedHits(t time.Duration) float64 {
	if t < 0 {
		return 0
	}

	x := min(t.Seconds(), p.RampUpDuration.Seconds())

	hitsInSlope := (p.a*math.Pow(x, 2))/2 + p.b*x

	if t < p.RampUpDuration {
		return hitsInSlope
	}

	linearDuration := t.Seconds() - p.RampUpDuration.Seconds()
	targetRate := p.TargetRate.hitsPerSec()

	return hitsInSlope + linearDuration*targetRate
}

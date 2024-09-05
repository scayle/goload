package pacer

import (
	"time"
)

// Pacer idea copied from https://github.com/tsenart/vegeta

// A Pacer defines the rate of hits during an Attack.
type Pacer interface {
	// Pace returns the duration an Runner should wait until
	// hitting the next Target, given an already elapsed duration and
	// completed hits.
	Pace(elapsed time.Duration, hits uint64) (wait time.Duration)

	// Rate returns a Pacer's instantaneous hit rate (per seconds)
	// at the given elapsed duration of an execution.
	Rate(elapsed time.Duration) float64
}

type Rate struct {
	Freq int           // Frequency (number of occurrences) per ...
	Per  time.Duration // Time unit, usually 1s
}

func (cp Rate) hitsPerSec() float64 {
	return (float64(cp.Freq) / float64(cp.Per)) * 1e9
}

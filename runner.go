package goload

import (
	"context"
	"github.com/HenriBeck/goload/pacer"
	"github.com/mroth/weightedrand/v2"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Runner struct {
	stopch          chan struct{}
	stopOnce        sync.Once
	workers         int
	maxWorkers      int
	weightOverrides map[string]int

	startedAt *time.Time
}

func NewRunner(loadTestOptions LoadTestOptions) *Runner {
	a := &Runner{
		stopch:          make(chan struct{}),
		stopOnce:        sync.Once{},
		workers:         loadTestOptions.initialWorkers,
		maxWorkers:      loadTestOptions.maxWorkers,
		weightOverrides: loadTestOptions.weightOverrides,
		startedAt:       nil,
	}

	return a
}

func (r *Runner) Run(ctx context.Context, exs []Executor, p pacer.Pacer, du time.Duration) <-chan *Result {
	var wg sync.WaitGroup

	workers := r.workers
	if workers > r.maxWorkers {
		workers = r.maxWorkers
	}

	chooser := r.getExecutorChooser(exs)

	now := time.Now()
	r.startedAt = &now

	results := make(chan *Result)
	ticks := make(chan struct{})
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go r.run(chooser, now, &wg, ticks, results)
	}

	go func() {
		<-ctx.Done()
		r.Stop()
	}()

	go func() {
		defer func() {
			close(ticks)
			wg.Wait()
			close(results)
			r.Stop()
		}()

		count := uint64(0)
		for {
			elapsed := time.Since(now)
			if du > 0 && elapsed > du {
				return
			}

			wait := p.Pace(elapsed, count)

			time.Sleep(wait)

			if workers < r.maxWorkers {
				select {
				case ticks <- struct{}{}:
					count++
					continue
				case <-r.stopch:
					return
				default:
					// all workers are blocked. start one more and try again
					workers++
					wg.Add(1)
					go r.run(chooser, now, &wg, ticks, results)
				}
			}

			select {
			case ticks <- struct{}{}:
				count++
			case <-r.stopch:
				return
			}
		}
	}()

	return results
}

func (r *Runner) getExecutorChooser(exs []Executor) *weightedrand.Chooser[Executor, int] {
	choices := make([]weightedrand.Choice[Executor, int], 0, len(exs))
	for _, ex := range exs {
		weight := ex.Options().Weight
		if override, ok := r.weightOverrides[ex.Name()]; ok {
			weight = override
		}
		choices = append(choices, weightedrand.NewChoice(ex, weight))
	}
	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create chooser")
	}
	return chooser
}

// Stop stops the current execution. The return value indicates whether this call
// has signalled the execution to stop (`true` for the first call) or whether it
// was a noop because it has been previously signalled to stop (`false` for any
// subsequent calls).
func (r *Runner) Stop() bool {
	select {
	case <-r.stopch:
		return false
	default:
		r.stopOnce.Do(func() { close(r.stopch) })
		return true
	}
}

func (r *Runner) run(chooser *weightedrand.Chooser[Executor, int], began time.Time, workers *sync.WaitGroup, ticks <-chan struct{}, results chan<- *Result) {
	defer workers.Done()

	for range ticks {
		results <- r.hit(chooser.Pick(), began)
	}
}

func (r *Runner) hit(ex Executor, began time.Time) *Result {
	res := Result{
		Timestamp: began.Add(time.Since(began)),
	}

	defer func() {
		res.Latency = time.Since(res.Timestamp)
	}()

	ctx := context.Background()
	if ex.Options().Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), ex.Options().Timeout)
		defer cancel()
	}

	resp := ex.Execute(ctx)

	res.Identifier = resp.Identifier
	res.AdditionalData = resp.AdditionalData
	res.Err = resp.Err

	return &res
}

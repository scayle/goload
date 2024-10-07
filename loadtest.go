package goload

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scayle/goload/pacer"
	ctx_utils "github.com/scayle/goload/utils/ctx"
	"math"
	"os"
	"time"
)

type resultHandler func(lt *LoadTest, result *Result)

type LoadTest struct {
	Pacer     pacer.Pacer
	Runner    *Runner
	Executors []Executor

	duration time.Duration

	resultHandlers   []resultHandler
	resultAggregator *resultAggregator
	reportInterval   time.Duration

	done chan struct{}
}

type LoadTestOptions struct {
	pacer           pacer.Pacer
	executors       []Executor
	duration        time.Duration
	initialWorkers  int
	maxWorkers      int
	resultHandlers  []resultHandler
	weightOverrides map[string]int
	reportInterval  time.Duration
	ctxModifier     func(ctx context.Context) context.Context
	defaultTimeout  time.Duration
}
type LoadTestOption func(*LoadTestOptions)

func RunLoadTest(opts ...LoadTestOption) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// set default options
	options, err := renderAndValidateOptions(opts)
	if err != nil {
		fmt.Printf("Invalid options: %v\n", err)
		os.Exit(1)
	}

	resultAggregator := newResultAggregator()
	options.resultHandlers = append(options.resultHandlers, resultAggregator.resultAggregationHandler)

	loadTest := &LoadTest{
		Pacer:            options.pacer,
		Runner:           NewRunner(options),
		Executors:        options.executors,
		duration:         options.duration,
		resultHandlers:   options.resultHandlers,
		resultAggregator: resultAggregator,
		reportInterval:   options.reportInterval,
		done:             make(chan struct{}),
	}

	ctx := ctx_utils.ContextWithInterrupt(context.Background())

	loadTest.Run(ctx)
}

func (lt *LoadTest) Run(ctx context.Context) {
	resultChan := lt.Runner.Run(ctx, lt.Executors, lt.Pacer, lt.duration)
	lt.runReporter(ctx)

	for result := range resultChan {
		for _, handler := range lt.resultHandlers {
			handler(lt, result)
		}
	}
	close(lt.done)
}

func (lt *LoadTest) runReporter(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(lt.reportInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Printf("expected pace: %.2f/s\n", lt.Pacer.Rate(time.Now().Sub(*lt.Runner.startedAt)))
				fmt.Printf("actual pace: %.2f/s\n", float64(lt.resultAggregator.rateCounter.Rate())/10)
				fmt.Printf("total hits: %d\n", lt.resultAggregator.total.Load())
				fmt.Printf("total failures: %d\n", lt.resultAggregator.failures.Load())
			}
		}
	}()
}

var defaultResultHandlers []resultHandler

func renderAndValidateOptions(opts []LoadTestOption) (LoadTestOptions, error) {
	options := LoadTestOptions{
		pacer:           nil,
		executors:       nil,
		duration:        0,
		initialWorkers:  10,
		maxWorkers:      math.MaxInt,
		resultHandlers:  defaultResultHandlers,
		weightOverrides: nil,
		reportInterval:  10 * time.Second,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if options.pacer == nil {
		return LoadTestOptions{}, fmt.Errorf("pacer is required")
	}
	if len(options.executors) == 0 {
		return LoadTestOptions{}, fmt.Errorf("should define at least one executor")
	}
	if options.initialWorkers == 0 || options.maxWorkers == 0 {
		return LoadTestOptions{}, fmt.Errorf("inital and max workers must be > 0")
	}

	return options, nil
}

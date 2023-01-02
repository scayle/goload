package goload

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type LoadTest struct {
	UI      *UI
	Options *LoadTestOptions
	Results *LoadTestResults

	done   chan bool
	ticker *time.Ticker
}

func RunLoadtest(configs ...LoadTestConfig) {
	options := &LoadTestOptions{}
	for _, config := range configs {
		config(options)
	}

	ui := NewUI(os.Stdout)

	ui.PrintStartMessage()

	loadtest := &LoadTest{
		Options: options,
		Results: NewResults(options.Endpoints),
		UI:      ui,

		done:   make(chan bool),
		ticker: initializeTicker(options, ui),
	}

	loadtest.WaitForLoadTestEnd()
	loadtest.ListenForAbort()

	loadtest.Run()

	ui.ReportResults(loadtest.Results)
}

func (loadtest *LoadTest) WaitForLoadTestEnd() {
	duration := loadtest.Options.LoadTestDuration
	if duration.Nanoseconds() == 0 {
		return
	}

	// Cancel the timer after the duration of the loadtest has elapsed
	go func() {
		time.Sleep(duration)
		loadtest.done <- true
		loadtest.ticker.Stop()
	}()
}

func (loadtest *LoadTest) ListenForAbort() {
	// Cancel the loadtest if the program is stopped using OS signals
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-ch

		fmt.Println()
		loadtest.UI.PrintAbortMessage()

		loadtest.done <- true
		loadtest.ticker.Stop()
	}()
}

func initializeTicker(options *LoadTestOptions, ui *UI) *time.Ticker {
	initialRPM := options.RPMStrategy.GetRPMForMinute(0)
	ticker := time.NewTicker(
		time.Minute / time.Duration(initialRPM),
	)

	ui.ReportInitialRPM(initialRPM)

	go func() {
		minute := int32(0)
		previousRPM := initialRPM
		t := time.NewTicker(time.Minute)

		for range t.C {
			minute++
			rpm := options.RPMStrategy.GetRPMForMinute(minute)

			if previousRPM == rpm {
				continue
			}

			if previousRPM < rpm {
				ui.ReportIncreaseInRPM(rpm)
			} else {
				ui.ReportDecreaseInRPM(rpm)
			}

			ticker.Reset(
				time.Minute / time.Duration(rpm),
			)
			previousRPM = rpm
		}
	}()

	return ticker
}

func (loadtest *LoadTest) Run() {
	endpointRandomizer := NewEndpointRandomizer(
		loadtest.Options.Endpoints,
	)
	g := new(sync.WaitGroup)

loop:
	for {
		select {
		case <-loadtest.done:
			break loop

		case <-loadtest.ticker.C:
			go func() {
				g.Add(1)
				defer g.Done()

				endpoint := endpointRandomizer.PickRandomEndpoint()

				ctx := context.Background()
				for _, fn := range loadtest.Options.ContextModifiers {
					ctx = fn(ctx)
				}

				if loadtest.Options.DefaultEndpointTimeout.Nanoseconds() > 0 {
					_ctx, cancel := context.WithTimeout(ctx, loadtest.Options.DefaultEndpointTimeout)
					defer cancel()

					ctx = _ctx
				}

				startTime := time.Now()
				err := endpoint.Execute(ctx)
				endTime := time.Now()

				loadtest.Results.SaveEndpointResult(
					endpoint,
					EndpointResult{
						Failed:   err != nil,
						Duration: endTime.Sub(startTime),
					},
				)
			}()
		}
	}

	// Wait until all requests have finished
	g.Wait()
}

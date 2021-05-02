package goload

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func RunLoadtest(
	configs ...LoadTestConfig,
) {
	rand.Seed(time.Now().UTC().UnixNano())

	options := &LoadTestOptions{}
	for _, config := range configs {
		config(options)
	}

	done := make(chan bool)
	timer := time.NewTicker(
		time.Minute / time.Duration(options.GetRequestsPerMinute()),
	)

	if options.LoadTestDuration.Nanoseconds() > 0 {
		// Cancel the timer after the duration of the loadtest has elapsed
		go func() {
			time.Sleep(options.LoadTestDuration)
			done <- true
			timer.Stop()
		}()
	}

	// Cancel the loadtest if the program is stopped using OS signals
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-ch

		fmt.Println("Cancelling timer")
		timer.Stop()
		done <- true
	}()

	endpointRandomizer := NewEndpointRandomizer(options.Endpoints)
	g := new(sync.WaitGroup)

loop:
	for {
		select {
		case <-done:
			break loop

		case <-timer.C:
			go func() {
				g.Add(1)
				defer g.Done()

				endpoint := endpointRandomizer.PickRandomEndpoint()
				ctx := context.Background()
				if timeout, ok := getTimeoutForEndpoint(endpoint, options); ok {
					_ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					ctx = _ctx
				}

				err := endpoint.Execute(ctx)
				if err != nil {
					fmt.Println(err)
				}
				// TODO: Record error
			}()
		}
	}

	g.Wait()
}

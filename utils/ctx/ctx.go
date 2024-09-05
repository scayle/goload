package ctx_utils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func ContextWithInterrupt(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)

	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGHUP)

		<-sigint

		fmt.Println("received interrupt signal: canceling context")

		cancel()
	}()

	return newCtx
}

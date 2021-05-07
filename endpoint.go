package goload

import (
	"context"
)

type Endpoint interface {
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	Execute(ctx context.Context) error

	GetRequestsPerMinute() int32

	Name() string
}

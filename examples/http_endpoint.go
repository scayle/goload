package main

import (
	"context"
	"fmt"
	"net/http"
)

type HTTPEndpoint struct{}

func (e *HTTPEndpoint) Execute(ctx context.Context) error {
	fmt.Println("Request")
	res, err := http.Get("http://test.k6.io")
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("non 200 status code")
	}

	return nil
}

func (e *HTTPEndpoint) GetRequestsPerMinute() int32 {
	return 16
}

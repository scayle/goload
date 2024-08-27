package goload_http

import (
	"context"
	"fmt"
	"github.com/HenriBeck/goload"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var basePath *string

type EndpointOption func(ep *endpoint)

func NewEndpoint(opts ...EndpointOption) goload.Executor {
	endpoint, err := renderAndValidateOptions(opts)
	if err != nil {
		fmt.Printf("Invalid Endpoint options: %v\n", err)
		os.Exit(1)
	}

	return endpoint
}

type endpoint struct {
	name    string
	weight  int
	timeout time.Duration

	client *http.Client

	urlFunc    func() *url.URL
	methodFunc func() string
	bodyFunc   func() io.Reader
	headerFunc func() http.Header

	validateResponse func(response *http.Response) error
}

func (e *endpoint) Execute(ctx context.Context) goload.ExecutionResponse {
	var body io.Reader
	if e.bodyFunc != nil {
		body = e.bodyFunc()
	}

	fullURL := e.urlFunc().String()
	req, err := http.NewRequestWithContext(ctx, e.methodFunc(), fullURL, body)
	if err != nil {
		return goload.ExecutionResponse{
			Identifier: e.name,
			Err:        err,
		}
	}

	if e.headerFunc != nil {
		req.Header = e.headerFunc()
	}

	res, err := e.client.Do(req)
	if err != nil {
		return goload.ExecutionResponse{
			Identifier: e.name,
			Err:        err,
		}
	}

	defer res.Body.Close()

	response := goload.ExecutionResponse{
		Identifier: e.name,
		Err:        nil,
		AdditionalData: map[string]string{
			"url": fullURL,
		},
	}

	if e.validateResponse != nil {
		if err := e.validateResponse(res); err != nil {
			response.Err = err
		}
	}

	return response
}

func (e *endpoint) Name() string {
	return e.name
}

func (e *endpoint) Options() *goload.ExecutorOptions {
	return &goload.ExecutorOptions{
		Weight:  e.weight,
		Timeout: e.timeout,
	}
}

func Status2xxResponseValidation(response *http.Response) error {
	if response.StatusCode < 200 && response.StatusCode >= 300 {
		return fmt.Errorf("non 2xx status code: %d", response.StatusCode)
	}
	return nil
}

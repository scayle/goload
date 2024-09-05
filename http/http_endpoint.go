package goload_http

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/scayle/goload"
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

	urlFunc    func() (*url.URL, error)
	methodFunc func() (string, error)
	bodyFunc   func() (io.Reader, error)
	headerFunc func() (http.Header, error)

	validateResponse func(response *http.Response) error
}

func (e *endpoint) Execute(ctx context.Context) goload.ExecutionResponse {
	response := goload.ExecutionResponse{
		Identifier: e.name,
	}

	var body io.Reader
	if e.bodyFunc != nil {
		var err error
		body, err = e.bodyFunc()
		if err != nil {
			response.Err = err
			log.Error().Err(err).Msg("failed to get body")
			return response
		}
	}

	targetURL, err := e.urlFunc()
	if err != nil {
		response.Err = err
		log.Error().Err(err).Msg("failed to get target URL")
		return response
	}
	targetURLStr := targetURL.String()

	method, err := e.methodFunc()
	if err != nil {
		response.Err = err
		log.Error().Err(err).Msg("failed to get method")
		return response
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURLStr, body)
	if err != nil {
		response.Err = err
		log.Error().Err(err).Msg("failed to create request")
		return response
	}

	if e.headerFunc != nil {
		headers, err := e.headerFunc()
		if err != nil {
			response.Err = err
			log.Error().Err(err).Msg("failed to get headers")
			return response
		}
		req.Header = headers
	}

	res, err := e.client.Do(req)
	if err != nil {
		response.Err = err
		log.Error().Err(err).Msg("failed to execute request")
		return response
	}

	defer res.Body.Close()

	response.AdditionalData = map[string]string{
		"url": targetURLStr,
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

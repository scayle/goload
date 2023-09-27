package goload_http

import (
	"math/rand"
	"net/http"
)

type ClientPool struct {
	clients []*http.Client
	rand    *rand.Rand
}

type HTTPTransportOption = func(request *http.Request)

// WithUserAgent sets the user agent on each of the request for the client pool.
func WithUserAgent(userAgent string) HTTPTransportOption {
	return func(request *http.Request) {
		if request.Header.Get("user-agent") == "" {
			request.Header.Set("user-agent", userAgent)
		}
	}
}

// WithHeader sets a custom header on each request.
func WithHeader(key string, value string) HTTPTransportOption {
	return func(request *http.Request) {
		request.Header.Add(key, value)
	}
}

// WithBasicAuth sets basic auth credentials on each request.
func WithBasicAuth(username string, password string) HTTPTransportOption {
	return func(request *http.Request) {
		request.SetBasicAuth(username, password)
	}
}

type transport struct {
	options []HTTPTransportOption

	innerTransport http.RoundTripper
}

func (transport *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	for _, option := range transport.options {
		option(request)
	}

	return transport.innerTransport.RoundTrip(request)
}

// NewHTTPClient creates a new http client for loadtesting.
//
// It allows to specify various request options which will be applied to all requests.
func NewHTTPClient(options ...HTTPTransportOption) *http.Client {
	return &http.Client{
		Transport: &transport{
			options:        options,
			innerTransport: http.DefaultTransport,
		},
	}
}

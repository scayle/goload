package goload_http

import "net/http"

var defaultClient = http.DefaultClient

func SetHTTPClient(client *http.Client) {
	defaultClient = client
}

type HTTPTransportOption = func(request *http.Request)

// WithClientUserAgent sets the user agent on each of the request
func WithClientUserAgent(userAgent string) HTTPTransportOption {
	return func(request *http.Request) {
		if request.Header.Get("user-agent") == "" {
			request.Header.Set("user-agent", userAgent)
		}
	}
}

// WithClientHeader sets a custom header on each request.
func WithClientHeader(key string, value string) HTTPTransportOption {
	return func(request *http.Request) {
		request.Header.Add(key, value)
	}
}

// WithClientBasicAuth sets basic auth credentials on each request.
func WithClientBasicAuth(username string, password string) HTTPTransportOption {
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

// NewClient creates a new http client for loadtesting with custom options.
//
// It allows to specify various request options which will be applied to all requests.
func NewClient(options ...HTTPTransportOption) *http.Client {
	return &http.Client{
		Transport: &transport{
			options:        options,
			innerTransport: http.DefaultTransport,
		},
	}
}

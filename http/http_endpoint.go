package goload_http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/HenriBeck/goload"
)

type Endpoint struct {
	options *EndpointOptions
}

type EndpointOptions struct {
	name              string
	requestsPerMinute int32

	client *http.Client

	getUrl func() url.URL
	method string
	body   io.Reader

	validateResponse func(response *http.Response) error
}

type HTTPEndpointConfig func(options *EndpointOptions)

// WithHTTPClient configures a static http client to be used for the loadtest endpoint.
func WithHTTPClient(client *http.Client) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.client = client
	}
}

// WithHTTPMethod sets the HTTP method used for the requests.
//
// By default, the `Endpoint` will use an `GET` request.
func WithHTTPMethod(method string) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.method = method
	}
}

// WithRequestsPerMinute configures the targeted requests per minute compared to other endpoints.
func WithRequestsPerMinute(rpm int32) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.requestsPerMinute = rpm
	}
}

// WithURL allows setting a static URL for the loadtest endpoint.
func WithURL(uri url.URL) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.getUrl = func() url.URL {
			return uri
		}
	}
}

// WithURLFunc allows for a dynamic creation of the URL per request made in the loadtest.
//
// An explicit endpoint name needs to be provided here for reporting purposes as an identifier.
func WithURLFunc(endpointName string, getUrl func() url.URL) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.name = endpointName
		options.getUrl = getUrl
	}
}

// WithURLString allows setting a static string URL for the loadtest endpoint.
func WithURLString(uri string) HTTPEndpointConfig {
	parsedUri, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	return func(options *EndpointOptions) {
		options.getUrl = func() url.URL {
			return *parsedUri
		}
	}
}

// WithValidateResponse allows to configure a custom check if the request should be counted as successful or not.
//
// By default, the request is successful if it returns a 2xx status code.
func WithValidateResponse(validate func(res *http.Response) error) HTTPEndpointConfig {
	return func(options *EndpointOptions) {
		options.validateResponse = validate
	}
}

// NewEndpoint creates a new HTTP based loadtest endpoint.
//
// To configure it you can use the functional options.
func NewEndpoint(configs ...HTTPEndpointConfig) goload.Endpoint {
	options := &EndpointOptions{
		method: http.MethodGet,
		body:   http.NoBody,
		client: http.DefaultClient,
		validateResponse: func(res *http.Response) error {
			if res.StatusCode < 200 || res.StatusCode > 299 {
				return fmt.Errorf("received non 200 status code from the server")
			}

			return nil
		},
	}
	for _, config := range configs {
		config(options)
	}

	return &Endpoint{
		options: options,
	}
}

func (endpoint *Endpoint) GetRequestsPerMinute() int32 {
	return endpoint.options.requestsPerMinute
}

func (endpoint *Endpoint) Name() string {
	// If an explicit endpoint name is provided we have dynamic URLs per request
	if endpoint.options.name != "" {
		return endpoint.options.name
	}

	// Otherwise this will resolve to a static URL which we will use as the endpoints name
	uri := endpoint.options.getUrl()

	return fmt.Sprintf(
		"%s %s",
		endpoint.options.method,
		uri.String(),
	)
}

func (endpoint *Endpoint) Execute(ctx context.Context) error {
	uri := endpoint.options.getUrl()

	req, err := http.NewRequestWithContext(
		ctx,
		endpoint.options.method,
		uri.String(),
		endpoint.options.body,
	)
	if err != nil {
		return err
	}

	res, err := endpoint.options.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return endpoint.options.validateResponse(res)
}

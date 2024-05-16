package goload_http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/HenriBeck/goload"
)

type EndpointOptions struct {
	name         string
	endpointOpts []goload.EndpointOption

	client *http.Client

	getUrl func() url.URL
	method string
	body   io.Reader
	header http.Header

	validateResponse func(response *http.Response) error
}

func (options *EndpointOptions) Name() string {
	// If an explicit endpoint name is provided we have dynamic URLs per request
	if options.name != "" {
		return options.name
	}

	// Otherwise this will resolve to a static URL which we will use as the endpoints name
	uri := options.getUrl()

	return fmt.Sprintf("%s %s", options.method, uri.String())
}

type HTTPEndpointOption func(options *EndpointOptions)

// WithHTTPClient configures a static http defaultClient to be used for the loadtest endpoint.
func WithHTTPClient(client *http.Client) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.client = client
	}
}

// WithHTTPMethod sets the HTTP method used for the requests.
//
// By default, the `Endpoint` will use an `GET` request.
func WithHTTPMethod(method string) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.method = method
	}
}

func WithHeader(key string, value string) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.header.Add(key, value)
	}
}

// WithRequestsPerMinute configures the targeted requests per minute compared to other endpoints.
func WithRequestsPerMinute(rpm int32) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.endpointOpts = append(options.endpointOpts, goload.WithRequestsPerMinute(rpm))
	}
}

// WithTimeout configures a specific timeout duration for the endpoint overriding the global config.
func WithTimeout(timeout time.Duration) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.endpointOpts = append(options.endpointOpts, goload.WithTimeout(timeout))
	}
}

// WithURL allows setting a static URL for the loadtest endpoint.
func WithURL(uri url.URL) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.getUrl = func() url.URL {
			return uri
		}
	}
}

// WithURLFunc allows for a dynamic creation of the URL per request made in the loadtest.
//
// An explicit endpoint name needs to be provided here for reporting purposes as an identifier.
func WithURLFunc(endpointName string, getUrl func() url.URL) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.name = endpointName
		options.getUrl = getUrl
	}
}

// WithURLString allows setting a static string URL for the loadtest endpoint.
func WithURLString(uri string) HTTPEndpointOption {
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
func WithValidateResponse(validate func(res *http.Response) error) HTTPEndpointOption {
	return func(options *EndpointOptions) {
		options.validateResponse = validate
	}
}

// NewEndpoint creates a new HTTP based loadtest endpoint.
//
// To configure it you can use the functional options.
func NewEndpoint(opts ...HTTPEndpointOption) goload.Endpoint {
	options := &EndpointOptions{
		method: http.MethodGet,
		body:   http.NoBody,
		client: defaultClient,
		validateResponse: func(res *http.Response) error {
			if res.StatusCode < 200 || res.StatusCode > 299 {
				return fmt.Errorf("received non 200 status code from the server")
			}

			return nil
		},
	}
	for _, config := range opts {
		config(options)
	}

	return goload.NewEndpoint(
		options.Name(),
		func(ctx context.Context) error {
			uri := options.getUrl()

			req, err := http.NewRequestWithContext(ctx, options.method, uri.String(), options.body)
			if err != nil {
				return err
			}

			req.Header = options.header

			res, err := options.client.Do(req)
			if err != nil {
				return err
			}

			defer res.Body.Close()

			return options.validateResponse(res)
		},
		options.endpointOpts...,
	)
}

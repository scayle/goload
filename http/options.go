package goload_http

import (
	"errors"
	"github.com/HenriBeck/goload"
	"github.com/HenriBeck/goload/http/url_builder"
	"io"
	"net/http"
	"net/url"
	"time"
)

func renderAndValidateOptions(opts []EndpointOption) (*endpoint, error) {
	endpoint := endpoint{
		name:    "",
		weight:  1,
		timeout: 0,
		client:  defaultClient,
		urlFunc: nil,
		methodFunc: func() string {
			return http.MethodGet
		},
		bodyFunc:         nil,
		headerFunc:       nil,
		validateResponse: nil,
	}

	for _, opt := range opts {
		opt(&endpoint)
	}

	if endpoint.urlFunc == nil {
		return nil, errors.New("urlFunc is required")
	}

	if endpoint.name == "" {
		endpoint.name = endpoint.urlFunc().Path
	}

	return &endpoint, nil
}

func WithName(name string) EndpointOption {
	return func(ep *endpoint) {
		ep.name = name
	}
}

func WithWeight(weight int) EndpointOption {
	return func(ep *endpoint) {
		ep.weight = weight
	}
}

func WithTimeout(timeout time.Duration) EndpointOption {
	return func(ep *endpoint) {
		ep.timeout = timeout
	}
}

func WithClient(client http.Client) EndpointOption {
	return func(ep *endpoint) {
		ep.client = &client
	}
}

func WithURL(rawURL string) EndpointOption {
	return func(ep *endpoint) {
		ep.urlFunc = func() *url.URL {
			if basePath != nil {
				var err error
				rawURL, err = url.JoinPath(*basePath, rawURL)
				if err != nil {
					panic(err)
				}
			}
			u, err := url.Parse(rawURL)
			if err != nil {
				panic(err)
			}
			return u
		}
	}
}

func WithURLFunc(urlFunc func() *url.URL) EndpointOption {
	return func(ep *endpoint) {
		ep.urlFunc = urlFunc
	}
}

func WithMethod(method string) EndpointOption {
	return func(ep *endpoint) {
		ep.methodFunc = func() string {
			return method
		}
	}
}

func WithMethodFunc(methodFunc func() string) EndpointOption {
	return func(ep *endpoint) {
		ep.methodFunc = methodFunc
	}
}

func WithBodyFunc(bodyFunc func() io.Reader) EndpointOption {
	return func(ep *endpoint) {
		ep.bodyFunc = bodyFunc
	}
}

func WithHeader(header http.Header) EndpointOption {
	return func(ep *endpoint) {
		ep.headerFunc = func() http.Header {
			return header
		}
	}
}

func WithHeaderFunc(headerFunc func() http.Header) EndpointOption {
	return func(ep *endpoint) {
		ep.headerFunc = headerFunc
	}
}

func WithValidateResponse(validationFunc func(response *http.Response) error) EndpointOption {
	return func(ep *endpoint) {
		ep.validateResponse = validationFunc
	}
}

func WithURLBuilder(opts ...url_builder.URLBuilderOption) EndpointOption {
	builder := url_builder.NewURLBuilder(opts)
	return func(ep *endpoint) {
		ep.urlFunc = func() *url.URL {
			return builder.Build(basePath)
		}
	}
}

func WithBasePath(path string) goload.LoadTestOption {
	return func(_ *goload.LoadTestOptions) {
		basePath = &path
	}
}

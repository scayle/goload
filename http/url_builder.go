package goload_http

import (
	"github.com/HenriBeck/goload/http/query_param"
	"github.com/HenriBeck/goload/utils/random"
	"net/url"
	"strings"
)

type URLBuilder struct {
	rawURL                  string
	urlParameterRandomizers []URLParameterRandomizer
	queryParams             []query_param.Builder
}

type URLBuilderOption func(*URLBuilder)

func NewURLBuilder(opts []URLBuilderOption) *URLBuilder {
	urlBuilder := URLBuilder{}

	for _, opt := range opts {
		opt(&urlBuilder)
	}

	if urlBuilder.rawURL == "" {
		panic("urlBuilder url must not be empty")
	}

	return &urlBuilder
}

func (builder *URLBuilder) Build() *url.URL {
	q := url.Values{}

	for _, param := range builder.queryParams {
		for key, values := range param.Build() {
			for _, value := range values {
				q.Add(key, value)
			}
		}
	}

	query := q.Encode()

	rawURL := builder.rawURL
	for _, u := range builder.urlParameterRandomizers {
		v := u.GetValue()
		rawURL = strings.Replace(rawURL, u.key, v, 1)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	u.RawQuery = query

	return u
}

type URLParameterRandomizer struct {
	key    string
	values []string
}

func (u *URLParameterRandomizer) GetValue() string {
	if len(u.values) == 0 {
		panic("Values cant be empty")
	}
	index := random.Number(0, int64(len(u.values)-1))
	return u.values[index]
}

package url_builder

import (
	"fmt"
	"github.com/HenriBeck/goload/utils/random"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

type URLBuilder struct {
	rawURL                  string
	urlParameterRandomizers []URLParameterRandomizer
	queryParams             []QueryParamBuilder
}

type URLBuilderOption func(*URLBuilder)

func NewURLBuilder(opts []URLBuilderOption) *URLBuilder {
	urlBuilder := URLBuilder{}

	for _, opt := range opts {
		opt(&urlBuilder)
	}

	if urlBuilder.rawURL == "" {
		log.Fatal().Msg("NewURLBuilder must include WithRawURL option")
	}

	return &urlBuilder
}

func (builder *URLBuilder) Build(basePath *string) (*url.URL, error) {
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
	if basePath != nil {
		var err error
		rawURL, err = url.JoinPath(*basePath, builder.rawURL)
		if err != nil {
			return nil, err
		}
	}

	for _, u := range builder.urlParameterRandomizers {
		v, err := u.GetValue()
		if err != nil {
			return nil, err
		}
		rawURL = strings.Replace(rawURL, u.key, v, 1)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	u.RawQuery = query

	return u, nil
}

type URLParameterRandomizer struct {
	key    string
	values []string
}

func (u *URLParameterRandomizer) GetValue() (string, error) {
	if len(u.values) == 0 {
		return "", fmt.Errorf("empty values for URLParameterRandomizer")
	}
	index := random.Number(0, int64(len(u.values)-1))
	return u.values[index], nil
}

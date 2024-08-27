package url_builder

func WithRawURL(rawURL string) URLBuilderOption {
	return func(builder *URLBuilder) {
		builder.rawURL = rawURL
	}
}

func WithQueryParams(queryParams ...QueryParamBuilder) URLBuilderOption {
	return func(builder *URLBuilder) {
		builder.queryParams = append(builder.queryParams, queryParams...)
	}
}

func WithURLParam(key string, values []string) URLBuilderOption {
	return func(builder *URLBuilder) {
		builder.urlParameterRandomizers = append(builder.urlParameterRandomizers, URLParameterRandomizer{key, values})
	}
}

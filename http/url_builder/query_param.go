package url_builder

import (
	"github.com/HenriBeck/goload/utils/random"
	"github.com/mroth/weightedrand/v2"
	"github.com/rs/zerolog/log"
	"net/url"
)

type QueryParamBuilder interface {
	Build() url.Values
}

type NameFn func() string
type ShouldBeUsedFn func() bool
type ValuesFn func() []string

type QueryParameter struct {
	Name         NameFn
	ShouldBeUsed ShouldBeUsedFn
	Value        ValuesFn
}

type QueryParameterOption func(param *QueryParameter)

func NewQueryParameter(opts ...QueryParameterOption) *QueryParameter {
	param := &QueryParameter{
		//set defaults
		ShouldBeUsed: UseAlways(),
	}
	for _, opt := range opts {
		opt(param)
	}
	if param.Name == nil || param.Value == nil {
		log.Fatal().Msg("NewQueryParameter must contain opts for name and value")
	}
	return param
}

func (s *QueryParameter) Build() url.Values {
	if !s.ShouldBeUsed() {
		return url.Values{}
	}
	values := s.Value()

	return url.Values{
		s.Name(): values,
	}
}

func UseAlways() ShouldBeUsedFn {
	return func() bool {
		return true
	}
}

type oneOfParam struct {
	params []QueryParamBuilder
}

func WithOneOfParam(params ...QueryParamBuilder) QueryParamBuilder {
	if len(params) == 0 {
		log.Fatal().Msg("WithOneOfParam must contain at least one parameter")
	}
	return &oneOfParam{params: params}
}

func (p *oneOfParam) Build() url.Values {
	index := random.Number(0, int64(len(p.params)-1))
	return p.params[index].Build()
}

type chanceParam struct {
	chance int
	param  QueryParamBuilder
	r      *weightedrand.Chooser[bool, int]
}

func NewParamWithUsageChange(chance int, param QueryParamBuilder) QueryParamBuilder {
	if chance > 100 || chance < 0 {
		log.Fatal().Msg("NewParamWithUsageChange chance value must be between 0 and 100")
	}
	r, err := weightedrand.NewChooser(
		weightedrand.NewChoice(true, chance),
		weightedrand.NewChoice(true, 100-chance),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create chooser")
	}

	return &chanceParam{chance: chance, param: param, r: r}
}

func (p *chanceParam) Build() url.Values {
	if p.r.Pick() {
		return p.param.Build()
	}
	return url.Values{}
}

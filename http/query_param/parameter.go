package query_param

import (
	"github.com/HenriBeck/goload/utils/random"
	"github.com/mroth/weightedrand/v2"
	"net/url"
)

type Builder interface {
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

type Opt func(param *QueryParameter)

func New(opts ...Opt) *QueryParameter {
	param := &QueryParameter{
		//set defaults
		ShouldBeUsed: UseAlways(),
	}
	for _, opt := range opts {
		opt(param)
	}
	if param.Name == nil || param.Value == nil {
		panic("query_param.New must contain opts for name and value")
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

func WithName(name string) Opt {
	return func(param *QueryParameter) {
		param.Name = func() string {
			return name
		}
	}
}

type pctParam struct {
	chance int
	param  Builder
	r      *weightedrand.Chooser[bool, int]
}

func NewPctParam(chance int, param Builder) Builder {
	if chance > 100 || chance < 0 {
		panic("chance value must be between 0 and 100")
	}
	r, err := weightedrand.NewChooser(
		weightedrand.NewChoice(true, chance),
		weightedrand.NewChoice(true, 100-chance),
	)
	if err != nil {
		panic(err)
	}

	return &pctParam{chance: chance, param: param, r: r}
}

func (p *pctParam) Build() url.Values {
	if p.r.Pick() {
		return p.param.Build()
	}
	return url.Values{}
}

type oneOfParam struct {
	params []Builder
}

func (o *oneOfParam) Build() url.Values {
	index := random.Number(0, int64(len(o.params)-1))
	return o.params[index].Build()
}

func NewOneOfParam(params ...Builder) Builder {
	if len(params) == 0 {
		panic("params cant be empty")
	}
	return &oneOfParam{params: params}
}

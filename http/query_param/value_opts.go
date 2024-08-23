package query_param

import (
	"github.com/HenriBeck/goload/utils/random"
	"github.com/mroth/weightedrand/v2"
	"strconv"
)

func WithValue(value string) Opt {
	return func(param *QueryParameter) {
		param.Value = func() []string {
			return []string{value}
		}
	}
}

func WithSampledValues(min int64, max int64, opts []string) Opt {
	sampler := random.NewSampler(opts)
	return func(param *QueryParameter) {
		n := random.Number(min, max)
		param.Value = func() []string {
			return sampler.Get(int(n))
		}
	}
}

type WeightedValueOpt struct {
	Value  string
	Weight int
}

func WithWeightedValue(opts ...WeightedValueOpt) Opt {
	if len(opts) == 0 {
		panic("WithWeightedValue opts can't be empty")
	}
	values := make([]weightedrand.Choice[string, int], 0, len(opts))
	for _, opt := range opts {
		values = append(values, weightedrand.NewChoice(opt.Value, opt.Weight))
	}

	r, err := weightedrand.NewChooser(values...)
	if err != nil {
		panic(err)
	}

	return func(param *QueryParameter) {
		param.Value = func() []string {
			return []string{r.Pick()}
		}
	}
}

func WithOneOfValue(values []string) Opt {
	if len(values) == 0 {
		panic("one off values needs at least one option")
	}
	return func(param *QueryParameter) {
		param.Value = func() []string {
			index := random.Number(0, int64(len(values)-1))
			return []string{values[index]}
		}
	}
}

func WithRandomNumberValue(min int64, max int64) Opt {
	return func(param *QueryParameter) {
		param.Value = func() []string {
			number := random.Number(min, max)
			return []string{strconv.FormatInt(number, 10)}
		}
	}
}

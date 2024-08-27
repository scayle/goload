package url_builder

import (
	"github.com/HenriBeck/goload/utils/random"
	"github.com/mroth/weightedrand/v2"
	"strconv"
)

func WithParamName(name string) QueryParameterOption {
	return func(param *QueryParameter) {
		param.Name = func() string {
			return name
		}
	}
}

func WithParamUsagePercentage(pct int) QueryParameterOption {
	if pct > 100 || pct < 0 {
		panic("WithParamUsagePercentage value must be between 0 and 100")
	}
	r, err := weightedrand.NewChooser(
		weightedrand.NewChoice(true, pct),
		weightedrand.NewChoice(true, 100-pct),
	)
	if err != nil {
		panic(err)
	}

	return func(param *QueryParameter) {
		param.ShouldBeUsed = func() bool {
			return r.Pick()
		}
	}
}

func WithParamValue(value string) QueryParameterOption {
	return func(param *QueryParameter) {
		param.Value = func() []string {
			return []string{value}
		}
	}
}

func WithSampledParamValues(min int64, max int64, opts []string) QueryParameterOption {
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

func WithWeightedParamValue(opts ...WeightedValueOpt) QueryParameterOption {
	if len(opts) == 0 {
		panic("WithWeightedParamValue opts can't be empty")
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

func WithOneOfParamValue(values []string) QueryParameterOption {
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

func WithRandomNumberParamValue(min int64, max int64) QueryParameterOption {
	return func(param *QueryParameter) {
		param.Value = func() []string {
			number := random.Number(min, max)
			return []string{strconv.FormatInt(number, 10)}
		}
	}
}

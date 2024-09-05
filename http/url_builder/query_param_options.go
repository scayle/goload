package url_builder

import (
	"github.com/mroth/weightedrand/v2"
	"github.com/rs/zerolog/log"
	"github.com/scayle/goload/utils/random"
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
		log.Fatal().Msg("WithParamUsagePercentage pct must be between 0 and 100")
	}
	r, err := weightedrand.NewChooser(
		weightedrand.NewChoice(true, pct),
		weightedrand.NewChoice(true, 100-pct),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create chooser")
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
		log.Fatal().Msg("WithWeightedParamValue must have at least one option")
	}
	values := make([]weightedrand.Choice[string, int], 0, len(opts))
	for _, opt := range opts {
		values = append(values, weightedrand.NewChoice(opt.Value, opt.Weight))
	}

	r, err := weightedrand.NewChooser(values...)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create chooser")
	}

	return func(param *QueryParameter) {
		param.Value = func() []string {
			return []string{r.Pick()}
		}
	}
}

func WithOneOfParamValue(values []string) QueryParameterOption {
	if len(values) == 0 {
		log.Fatal().Msg("WithOneOfParamValue must have at least one value")
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

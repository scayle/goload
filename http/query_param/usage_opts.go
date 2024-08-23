package query_param

import "github.com/mroth/weightedrand/v2"

func WithUsagePct(pct int) Opt {
	if pct > 100 || pct < 0 {
		panic("WithUsagePct value must be between 0 and 100")
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

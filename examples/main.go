package main

import (
	"fmt"
	"github.com/scayle/goload"
	goload_http "github.com/scayle/goload/http"
	"github.com/scayle/goload/http/url_builder"
	"github.com/scayle/goload/pacer"
	"time"
)

func main() {
	goload.RunLoadTest(
		goload.WithDuration(5*time.Minute),
		//goload.WithLinearRampUpPacer(pacer.Rate{Freq: 30, Per: time.Minute}, pacer.Rate{Freq: 2, Per: time.Second}, 1*time.Minute),
		goload.WithConstantPacer(pacer.Rate{Freq: 2, Per: time.Second}),
		goload_http.WithBasePath("http://test.k6.io"),
		goload.WithExecutors(
			goload_http.NewEndpoint(
				goload_http.WithName("test"),
				goload_http.WithURL("/"),
				goload_http.WithValidateResponse(goload_http.Status2xxResponseValidation),
			),
			goload_http.NewEndpoint(
				goload_http.WithName("pi"),
				goload_http.WithURLBuilder(
					url_builder.WithRawURL("/pi.php"),
					url_builder.WithQueryParams(
						url_builder.NewQueryParameter(
							url_builder.WithParamName("decimals"),
							url_builder.WithRandomNumberParamValue(1, 20),
						),
					),
				),
				goload_http.WithValidateResponse(goload_http.Status2xxResponseValidation),
			),
		),
		goload.WithAdditionalResultHandler(func(_ *goload.LoadTest, result *goload.Result) {
			fmt.Printf("result: %+v\n", result)
		}),
		goload.WithWeightOverrides(map[string]int{
			"test": 1,
			"pi":   2,
		}),
	)
}

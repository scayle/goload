package main

import (
	"fmt"
	"github.com/HenriBeck/goload"
	goload_http "github.com/HenriBeck/goload/http"
	"github.com/HenriBeck/goload/http/query_param"
	"github.com/HenriBeck/goload/pacer"
	"time"
)

func main() {
	goload.RunLoadTest(
		goload.WithDuration(5*time.Minute),
		//goload.WithLinearRampUpPacer(pacer.Rate{Freq: 30, Per: time.Minute}, pacer.Rate{Freq: 2, Per: time.Second}, 1*time.Minute),
		goload.WithConstantPacer(pacer.Rate{Freq: 2, Per: time.Second}),
		goload.WithExecutors(
			goload_http.NewEndpoint(
				goload_http.WithName("test"),
				goload_http.WithURL("http://test.k6.io"),
				goload_http.WithValidateResponse(goload_http.Status2xxResponseValidation),
			),
			goload_http.NewEndpoint(
				goload_http.WithName("pi"),
				goload_http.WithURLBuilder(
					goload_http.WithRawURL("https://test.k6.io/pi.php"),
					goload_http.WithQueryParams(
						query_param.New(
							query_param.WithName("decimals"),
							query_param.WithRandomNumberValue(1, 20),
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

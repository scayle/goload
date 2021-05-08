package main

import (
	"time"

	"github.com/HenriBeck/goload"
	goload_http "github.com/HenriBeck/goload/http"
)

func main() {
	goload.RunLoadtest(
		goload.WithDuration(5*time.Minute),
		goload.WithEndpoint(
			goload_http.NewHTTPEndpoint(
				goload_http.WithURLString("http://test.k6.io"),
				goload_http.WithRequestsPerMinute(15),
			),
		),
		goload.WithEndpoint(
			goload_http.NewHTTPEndpoint(
				goload_http.WithURLString("http://test.k6.io/news.php"),
				goload_http.WithRequestsPerMinute(10),
			),
		),
		goload.WithRampUpRPM(
			[]goload.Step{
				{Minute: 1, RPM: 10},
				{Minute: 2, RPM: 12},
				{Minute: 3, RPM: 6},
			},
		),
	)
}

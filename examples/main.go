package main

import (
	"time"

	"github.com/HenriBeck/goload"
)

func main() {
	goload.RunLoadtest(
		goload.WithDuration(10*time.Minute),
		goload.WithEndpoint(&HTTPEndpoint{}),
		goload.WithRampUpRPM(
			[]goload.Step{
				{Minute: 1, RPM: 10},
				{Minute: 2, RPM: 12},
				{Minute: 3, RPM: 6},
			},
		),
	)
}

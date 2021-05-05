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
				{Minute: 2, RPM: 10},
				{Minute: 3, RPM: 12},
				{Minute: 6, RPM: 16},
				{Minute: 8, RPM: 20},
			},
		),
	)
}

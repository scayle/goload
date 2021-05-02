package main

import (
	"time"

	"github.com/HenriBeck/goload"
)

func main() {
	goload.RunLoadtest(
		goload.WithDuration(2*time.Minute),
		goload.WithEndpoint(&HTTPEndpoint{}),
	)
}

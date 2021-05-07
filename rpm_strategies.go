package goload

import (
	"math"
	"sort"
)

type RPMStrategy interface {
	GetRPMForMinute(minute int32) int32
}

type StaticRPMStrategy struct {
	rpm int32
}

func NewStaticRPMStrategy(rpm int32) RPMStrategy {
	return &StaticRPMStrategy{
		rpm: rpm,
	}
}

func (strategy *StaticRPMStrategy) GetRPMForMinute(minute int32) int32 {
	return strategy.rpm
}

type RampUpRPMStrategy struct {
	steps []Step
}

func NewRampUpRPMStrategy(steps []Step) RPMStrategy {
	isSorted := sort.SliceIsSorted(steps, func(i, j int) bool {
		return steps[i].Minute < steps[j].Minute
	})

	if !isSorted {
		panic("config isn't sorted")
	}

	return &RampUpRPMStrategy{
		steps: steps,
	}
}

type Step struct {
	Minute int32
	RPM    int32
}

func (strategy *RampUpRPMStrategy) GetRPMForMinute(minute int32) int32 {
	// In case there is a direct match we just return the RPM count for it
	for _, step := range strategy.steps {
		if step.Minute == minute {
			return step.RPM
		}
	}

	// If the `minute` is before our first step we just return the starting RPM
	firstStep := strategy.steps[0]
	if firstStep.Minute > minute {
		return firstStep.RPM
	}

	// Otherwise we look for the steps which encloses the current `minute`
	for index := 0; index < len(strategy.steps)-1; index++ {
		nextIndex := index + 1
		step := strategy.steps[index]
		nextStep := strategy.steps[nextIndex]

		// Check if the `step` is before the `minute` and if the `nextStep` is after the `minute`
		if step.Minute < minute && minute < nextStep.Minute {
			// We need to calculate the correct interpolation between the two steps and our wanted minute.
			//
			// As an example:
			// Step Minute 1 has 2 RPM
			// We are currently at Minute 2
			// Step for Minute 3 has 8 RPM
			//
			// In this case we need to take the difference between our two steps,
			// and then interpolate it in relation to the `minute`.
			// We then calculate how much RPM we need to add on top of our previous `step`.
			//
			// The example above would result in 5 RPM. (2 + 6 * 1 / 2)
			//
			// This also works for the reverse in case this is used with a ramp down steps.
			rpmDiffBetweenPoints := float64(nextStep.RPM) - float64(step.RPM)
			timeDiffBetweenPoints := float64(nextStep.Minute) - float64(step.Minute)
			timeDiffToPreviousPoint := float64(minute) - float64(step.Minute)

			return int32(
				float64(step.RPM) + math.Floor(rpmDiffBetweenPoints*timeDiffToPreviousPoint/timeDiffBetweenPoints),
			)
		}
	}

	// If we don't find steps that encloses the current `minute`,
	// it means that the `minute` is larger than our last step
	// so we just return the RPM for it.
	lastPoint := strategy.steps[len(strategy.steps)-1]
	return lastPoint.RPM
}

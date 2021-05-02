package goload

import (
	"sort"
	"time"
)

type LoadTestOptions struct {
	Endpoints []Endpoint

	DefaultEndpointTimeout time.Duration

	LoadTestDuration time.Duration

	RPMStrategy RPMStrategy
}

type RPMStrategy interface {
	GetStartingRPM() uint
	GetRPMForMinute(minute uint) uint
}

type LoadTestConfig func(options *LoadTestOptions)

func WithEndpoint(endpoint Endpoint) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.Endpoints = append(options.Endpoints, endpoint)
	}
}

func WithDefaultEndpointTimeout(timeout time.Duration) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.DefaultEndpointTimeout = timeout
	}
}

func WithDuration(loadTestDuration time.Duration) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.LoadTestDuration = loadTestDuration
	}
}

type StaticRPMStrategy struct {
	rpm uint
}

func (strategy *StaticRPMStrategy) GetStartingRPM() uint {
	return strategy.rpm
}

func (strategy *StaticRPMStrategy) GetRPMForMinute(minute uint) uint {
	return strategy.rpm
}

func WithStaticRPM(rpm uint) LoadTestConfig {
	return func(options *LoadTestOptions) {
		options.RPMStrategy = &StaticRPMStrategy{
			rpm: rpm,
		}
	}
}

type RampUpRPMStrategy struct {
	points []RampUpPoints
}

type RampUpPoints struct {
	Minute uint
	RPM    uint
}

func (strategy *RampUpRPMStrategy) GetStartingRPM() uint {
	return strategy.points[0].RPM
}

func (strategy *RampUpRPMStrategy) GetRPMForMinute(minute uint) uint {
	// In case there is a direct match
	for _, point := range strategy.points {
		if point.Minute == minute {
			return point.RPM
		}
	}

	firstPoint := strategy.points[0]
	if firstPoint.Minute > minute {
		return firstPoint.RPM
	}

	for index := 0; index < len(strategy.points)-1; index++ {
		nextIndex := index + 1
		point := strategy.points[index]
		nextPoint := strategy.points[nextIndex]

		if point.Minute < minute && minute < nextPoint.Minute {
			rpmDiffBetweenPoints := nextPoint.RPM - point.RPM
			timeDiffBetweenPoints := nextPoint.Minute - point.Minute
			timeDiffToPreviousPoint := point.Minute - minute

			return point.RPM + (rpmDiffBetweenPoints * timeDiffToPreviousPoint / timeDiffBetweenPoints)
		}
	}

	lastPoint := strategy.points[len(strategy.points)-1]
	return lastPoint.RPM
}

func WithRampUp(points []RampUpPoints) LoadTestConfig {
	isSorted := sort.SliceIsSorted(points, func(i, j int) bool {
		return points[i].Minute < points[j].Minute
	})

	if !isSorted {
		panic("config isn't sorted")
	}

	return func(options *LoadTestOptions) {
		options.RPMStrategy = &RampUpRPMStrategy{
			points: points,
		}
	}
}

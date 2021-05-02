package goload

import (
	"math"
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
	GetStartingRPM() int32
	GetRPMForMinute(minute int32) int32
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
	rpm int32
}

func (strategy *StaticRPMStrategy) GetStartingRPM() int32 {
	return strategy.rpm
}

func (strategy *StaticRPMStrategy) GetRPMForMinute(minute int32) int32 {
	return strategy.rpm
}

func WithStaticRPM(rpm int32) LoadTestConfig {
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
	Minute int32
	RPM    int32
}

func (strategy *RampUpRPMStrategy) GetStartingRPM() int32 {
	return strategy.points[0].RPM
}

func (strategy *RampUpRPMStrategy) GetRPMForMinute(minute int32) int32 {
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
			rpmDiffBetweenPoints := float64(nextPoint.RPM) - float64(point.RPM)
			timeDiffBetweenPoints := float64(nextPoint.Minute) - float64(point.Minute)
			timeDiffToPreviousPoint := float64(minute) - float64(point.Minute)

			return int32(
				float64(point.RPM) + math.Floor(rpmDiffBetweenPoints*timeDiffToPreviousPoint/timeDiffBetweenPoints),
			)
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

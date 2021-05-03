package goload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRampUpRPMStrategy(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		strategy := &RampUpRPMStrategy{
			points: []RampUpPoints{
				{Minute: 2, RPM: 10},
				{Minute: 3, RPM: 12},
				{Minute: 6, RPM: 16},
				{Minute: 8, RPM: 20},
			},
		}

		assert.Equal(t, int32(10), strategy.GetRPMForMinute(1))
		assert.Equal(t, int32(10), strategy.GetRPMForMinute(2))
		assert.Equal(t, int32(12), strategy.GetRPMForMinute(3))
		assert.Equal(t, int32(13), strategy.GetRPMForMinute(4))
		assert.Equal(t, int32(14), strategy.GetRPMForMinute(5))
		assert.Equal(t, int32(16), strategy.GetRPMForMinute(6))
		assert.Equal(t, int32(18), strategy.GetRPMForMinute(7))
		assert.Equal(t, int32(20), strategy.GetRPMForMinute(8))
		assert.Equal(t, int32(20), strategy.GetRPMForMinute(9))
		assert.Equal(t, int32(20), strategy.GetRPMForMinute(10))
	})
}

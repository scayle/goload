package goload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticRPMStrategy(t *testing.T) {
	t.Run("GetRPMForMinute should return the static RPM for any minute", func(t *testing.T) {
		strategy := NewStaticRPMStrategy(50)

		assert.Equal(
			t,
			int32(50),
			strategy.GetRPMForMinute(0),
		)
		assert.Equal(
			t,
			int32(50),
			strategy.GetRPMForMinute(1),
		)
		assert.Equal(
			t,
			int32(50),
			strategy.GetRPMForMinute(5),
		)
	})
}

func TestRampUpRPMStrategy(t *testing.T) {
	t.Run("should panic if the points are not sorted properly", func(t *testing.T) {
		assert.PanicsWithValue(
			t,
			"config isn't sorted",
			func() {
				NewRampUpRPMStrategy(
					[]Step{
						{Minute: 12},
						{Minute: 5},
					},
				)
			},
		)
	})
}

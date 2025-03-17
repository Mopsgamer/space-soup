package soup

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRadiansFromRich(t *testing.T) {
	assert := assert.New(t)

	// 34°10'16” == 0,596398 rad
	d := RadiansFromRich(34, 10, 16)

	assert.InDelta(d, 596e-3, 1e-3)
}

func TestRichFromRadians(t *testing.T) {
	assert := assert.New(t)

	// 0,596398 rad == 34°10'16”
	d, m, s := RichFromRadians(RadiansFromRich(34, 10, 16))

	assert.InDelta(d, 34, 1)
	assert.InDelta(m, 10, 1)
	assert.InDelta(s, 16, 1)
}

func TestLoopNumber(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0, LoopNumber(0, 0, 360))
	assert.Equal(10, LoopNumber(10, 0, 360))
	assert.Equal(0, LoopNumber(360, 0, 360))
	assert.Equal(10, LoopNumber(370, 0, 360))
	assert.Equal(350, LoopNumber(-10, 0, 360))

	assert.Equal(-3, LoopNumber(-1, -3, -1))
	assert.Equal(-2, LoopNumber(-2, -3, -1))
	assert.Equal(-3, LoopNumber(-3, -3, -1))

	assert.Equal(-math.Pi/2, LoopNumber(-math.Pi/2, -math.Pi/2, math.Pi/2))
	assert.Equal(-math.Pi/2, LoopNumber(math.Pi/2, -math.Pi/2, math.Pi/2))
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Float64("12"), 12.)
	assert.Equal(Float64("12.64"), 12.64)
	assert.Equal(Float64("12,64"), 12.64)
}

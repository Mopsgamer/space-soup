package soup

import (
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

	assert.Equal(LoopNumber(0, 0, 360), 0)
	assert.Equal(LoopNumber(10, 0, 360), 10)
	assert.Equal(LoopNumber(360, 0, 360), 0)
	assert.Equal(LoopNumber(370, 0, 360), 10)
	assert.Equal(LoopNumber(-10, 0, 360), 350)
}

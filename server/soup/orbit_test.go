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

func TestOrbit3(t *testing.T) {
	assert := assert.New(t)

	date, err := ParseDate("1972-01-25T06:07")
	if !assert.NoError(err) {
		return // error
	}

	movement := NewMovement(Input{
		Tau1:  -12.7572,
		Tau2:  -17.5536,
		V_avg: Average([]float64{33.858, 33.832, 33.965}),
		Date:  date,
	})

	assert.InDelta(59.827, DegreesFromRadians(movement.A), 1e-3)     // FIXME: Fails: 59.8242
	assert.InDelta(27.004, DegreesFromRadians(movement.Z_avg), 1e-3) // FIXME: Fails: 26.3362
}

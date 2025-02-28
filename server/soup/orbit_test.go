package soup

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRadiansFromRich(t *testing.T) {
	assert := assert.New(t)

	// 34°10'16” == 0,596398 rad
	assert.Equal(int(math.Floor(phi1*1000)), 596)
}

func TestOrbit3(t *testing.T) {
	assert := assert.New(t)

	date, err := ParseDate("1972-05-13T06:07")
	if !assert.NoError(err) {
		return // error
	}

	movement := NewMovement(Input{
		Tau1:  -12.7572,
		Tau2:  -17.5536,
		V_avg: Average([]float64{33.858, 33.832, 33.965}),
		Date:  date,
	})

	assert.Equal(59.827, DegreesFromRadians(movement.A)) // FIXME: test fails
}

func TestOrbit4(t *testing.T) {
	assert := assert.New(t)

	date, err := ParseDate("1972-05-13T06:27")
	if !assert.NoError(err) {
		return // error
	}

	movement := NewMovement(Input{
		Tau1:  -12.5536,
		Tau2:  -0.3927,
		V_avg: Average([]float64{56.36, 60.908, 55.398}),
		Date:  date,
	})

	assert.Equal(21.169, DegreesFromRadians(movement.A)) // FIXME: test fails
}

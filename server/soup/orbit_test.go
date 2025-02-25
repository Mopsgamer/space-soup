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
		Tau1:  -12.76,
		Tau2:  -17.55,
		V_avg: Average([]float64{33.86, 33.83, 33.97}),
		Date:  date,
	})

	assert.Equal(59.827, movement.A) // FIXME: test fails
}

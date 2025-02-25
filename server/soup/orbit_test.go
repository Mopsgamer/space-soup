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

func TestOrbit(t *testing.T) {
	assert := assert.New(t)

	date, err := ParseDate("1972-05-13T00:37")
	if !assert.NoError(err) {
		return // error
	}

	movement := NewMeteoroidMovement(MeteoroidMovementInput{
		Dist:  17299,
		Tau1:  471889,
		Tau2:  213476,
		V_avg: Average([]float64{18177, 999999, 11987}),
		Date:  date,
	})

	assert.Equal(33419, movement.A) // TODO: test fails
}

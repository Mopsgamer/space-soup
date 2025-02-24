package soup

import (
	"math"
	"reflect"
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

	// TODO: write tests
	calc := NewMeteoroidMovement(MeteoroidMovementInput{})

	assert.True(reflect.DeepEqual(*calc, MeteoroidMovement{}))
}

package soup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrbit(t *testing.T) {
	assert := assert.New(t)

	// TODO: write tests
	calc, err := NewMeteoroidMovement(MeteoroidMovementInput{})

	if assert.NoError(err) {
		assert.Equal(calc, &MeteoroidMovement{})
	}
}

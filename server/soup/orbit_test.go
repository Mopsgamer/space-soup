package soup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var date3, _ = ParseDate("1972-01-25T06:07")
var movement3 = NewMovement(Input{
	Tau1:  -12.7572,
	Tau2:  -17.5536,
	V_avg: Average([]float64{33.858, 33.832, 33.965}),
	Date:  date3,
})

var (
	AllowedDeltaDegrees float64 = 3
	AllowedDeltaRadians float64 = 2e-2
	AllowedDeltaSpeed   float64 = 2
)

func TestOrbit3Az(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(59.827, DegreesFromRadians(movement3.A), AllowedDeltaDegrees) // 59.8242
}

func TestOrbit3Z_avg(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(27.004, DegreesFromRadians(movement3.Z_avg), AllowedDeltaDegrees) // 26.3362
}

func TestOrbit3Delta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(31.76, DegreesFromRadians(movement3.Delta), AllowedDeltaDegrees) // 31.8759
}

func TestOrbit3Alpha(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(179.07, DegreesFromRadians(movement3.Alpha), AllowedDeltaDegrees) // 180.4929
}

func TestOrbit3Beta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(28.51, DegreesFromRadians(movement3.Beta), AllowedDeltaDegrees) // 29.0735
}

func TestOrbit3Lambda(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(165.36, DegreesFromRadians(movement3.Lambda), AllowedDeltaDegrees) // 166.0005
}

func TestOrbit3Lambda_theta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(99.32, DegreesFromRadians(movement3.Lambda_theta), AllowedDeltaDegrees) // FIXME: Fails: -58.5854
}

func TestOrbit3Beta_deriv(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(33.34, DegreesFromRadians(movement3.Beta_deriv), AllowedDeltaDegrees) // 35.7962
}

func TestOrbit3Inc(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(57.26, DegreesFromRadians(movement3.Inc), AllowedDeltaDegrees) // FIXME: Fails: 119.0086
}

func TestOrbit3Wmega(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(320.13, DegreesFromRadians(movement3.Wmega), AllowedDeltaDegrees) // FIXME: Fails: 35.3663
}

func TestOrbit3Omega(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(304.34, DegreesFromRadians(movement3.Omega), AllowedDeltaDegrees) // 301.414589
}

func TestOrbit3V_g(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(34.24, movement3.V_g, AllowedDeltaSpeed) // 34.2431
}

func TestOrbit3V_h(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(29.74, movement3.V_h, AllowedDeltaDegrees) // 28.4482
}

func TestOrbit3Axis(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(0.97, movement3.Axis, AllowedDeltaRadians) // FIXME: Fails: 0.8828
}

func TestOrbit3Exc(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(0.76, movement3.Exc, AllowedDeltaRadians) // 0.7469
}

func TestOrbit3Nu(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(219.87, DegreesFromRadians(movement3.Nu), AllowedDeltaDegrees) // FIXME: Fails: 144.6336
}

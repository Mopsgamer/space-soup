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

var date3, _ = ParseDate("1972-01-25T06:07")
var movement3 = NewMovement(Input{
	Tau1:  -12.7572,
	Tau2:  -17.5536,
	V_avg: Average([]float64{33.858, 33.832, 33.965}),
	Date:  date3,
})

func TestOrbit3Az(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(59.827, DegreesFromRadians(movement3.A), 1) // 59.8242
}

func TestOrbit3Z_avg(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(27.004, DegreesFromRadians(movement3.Z_avg), 1) // 26.3362
}

func TestOrbit3Delta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(31.76, DegreesFromRadians(movement3.Delta), 1) // 31.8759
}

func TestOrbit3Alpha(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(179.07, DegreesFromRadians(movement3.Alpha), 1e-2) // FIXME: Fails: 92.4567
}

func TestOrbit3Beta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(28.51, DegreesFromRadians(movement3.Beta), 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Lambda(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(165.36, DegreesFromRadians(movement3.Lambda), 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Lambda_theta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(99.32, DegreesFromRadians(movement3.Lambda_theta), 1e-2) // FIXME: Fails: -1.0181
}

func TestOrbit3Beta_deriv(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(33.34, DegreesFromRadians(movement3.Beta_deriv), 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Inc(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(57.26, DegreesFromRadians(movement3.Inc), 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Wmega(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(320.13, DegreesFromRadians(movement3.Wmega), 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Omega(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(304.34, DegreesFromRadians(movement3.Omega), 1e-2) // FIXME: Fails: -1.0181
}

func TestOrbit3V_g(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(34.24, movement3.V_g, 1e-2) // 34.2431
}

func TestOrbit3V_h(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(29.74, movement3.V_h, 1e-3) // FIXME: Fails: NaN
}

func TestOrbit3Axis(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(0.97, movement3.Axis, 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Exc(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(0.76, movement3.Exc, 1e-2) // FIXME: Fails: NaN
}

func TestOrbit3Nu(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(219.87, movement3.Nu, 1e-2) // FIXME: Fails: NaN
}

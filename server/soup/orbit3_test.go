package soup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var date3, _ = ParseDateJSON("1972-01-25T06:07")
var movement3 = NewMovement(Input{
	Tau1:  -12.7572,
	Tau2:  -17.5536,
	V_avg: Average([]float64{33.858, 33.832, 33.965}),
	Date:  date3,
})

var (
	allowedDeltaDegrees float64 = 4
	allowedDeltaRadians float64 = 2e-2
	allowedDeltaSpeed   float64 = 2
)

func TestOrbit3(t *testing.T) {
	CheckOrbit(t, &Movement[float64]{
		Lambda_apex:  RadiansFromDegrees(213.99),
		A:            RadiansFromDegrees(59.827),
		Z_avg:        RadiansFromDegrees(27.004),
		Delta:        RadiansFromDegrees(31.76),
		Alpha:        RadiansFromDegrees(179.07),
		Beta:         RadiansFromDegrees(28.51),
		Lambda:       RadiansFromDegrees(165.36),
		Lambda_theta: RadiansFromDegrees(99.32),
		Beta_deriv:   RadiansFromDegrees(33.34),
		Inc:          RadiansFromDegrees(57.26),
		Wmega:        RadiansFromDegrees(320.13),
		Omega:        RadiansFromDegrees(304.34),
		V_g:          34.24,
		V_h:          29.74,
		Axis:         0.97,
		Exc:          0.76,
		Nu:           RadiansFromDegrees(219.87),
	}, movement3)
}

func TestOrbit3Lambda_apex(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(213.99, DegreesFromRadians(movement3.Lambda_apex), allowedDeltaDegrees)
}

func TestOrbit3A(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(59.827, DegreesFromRadians(movement3.A), allowedDeltaDegrees)
}

func TestOrbit3Z_avg(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(27.004, DegreesFromRadians(movement3.Z_avg), allowedDeltaDegrees)
}

func TestOrbit3Delta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(31.76, DegreesFromRadians(movement3.Delta), allowedDeltaDegrees)
}

func TestOrbit3Alpha(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(179.07, DegreesFromRadians(movement3.Alpha), allowedDeltaDegrees)
}

func TestOrbit3Beta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(28.51, DegreesFromRadians(movement3.Beta), allowedDeltaDegrees)
}

func TestOrbit3Lambda(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(165.36, DegreesFromRadians(movement3.Lambda), allowedDeltaDegrees)
}

func TestOrbit3Lambda_theta(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(99.32, DegreesFromRadians(movement3.Lambda_theta), allowedDeltaDegrees) // FIXME: Fails: -58.5854
}

func TestOrbit3Beta_deriv(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(33.34, DegreesFromRadians(movement3.Beta_deriv), allowedDeltaDegrees)
}

func TestOrbit3Inc(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(57.26, DegreesFromRadians(movement3.Inc), allowedDeltaDegrees) // FIXME: Fails: 118.8793
}

func TestOrbit3Wmega(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(320.13, DegreesFromRadians(movement3.Wmega), allowedDeltaDegrees) // FIXME: Fails: 35.2268
}

func TestOrbit3Omega(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(304.34, DegreesFromRadians(movement3.Omega), allowedDeltaDegrees)
}

func TestOrbit3V_g(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(34.24, movement3.V_g, allowedDeltaSpeed)
}

func TestOrbit3V_h(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(29.74, movement3.V_h, allowedDeltaDegrees)
}

func TestOrbit3Axis(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(0.97, movement3.Axis, allowedDeltaRadians) // FIXME: Fails: 0.88028544
}

func TestOrbit3Exc(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(0.76, movement3.Exc, allowedDeltaRadians)
}

func TestOrbit3Nu(t *testing.T) {
	assert := assert.New(t)
	assert.InDelta(219.87, DegreesFromRadians(movement3.Nu), allowedDeltaDegrees) // FIXME: Fails: 144.7731
}

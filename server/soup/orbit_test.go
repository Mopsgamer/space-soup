package soup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CheckOrbit(t *testing.T, excpected, actual *Movement[float64]) {
	assert := assert.New(t)
	assert.InDelta(DegreesFromRadians(excpected.Lambda_apex), DegreesFromRadians(actual.Lambda_apex), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.A), DegreesFromRadians(actual.A), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Z_avg), DegreesFromRadians(actual.Z_avg), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Delta), DegreesFromRadians(actual.Delta), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Alpha), DegreesFromRadians(actual.Alpha), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Beta), DegreesFromRadians(actual.Beta), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Lambda), DegreesFromRadians(actual.Lambda), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Lambda_theta), DegreesFromRadians(actual.Lambda_theta), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Beta_deriv), DegreesFromRadians(actual.Beta_deriv), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Inc), DegreesFromRadians(actual.Inc), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Wmega), DegreesFromRadians(actual.Wmega), allowedDeltaDegrees)
	assert.InDelta(DegreesFromRadians(excpected.Omega), DegreesFromRadians(actual.Omega), allowedDeltaDegrees)
	assert.InDelta(excpected.V_g, actual.V_g, allowedDeltaSpeed)
	assert.InDelta(excpected.V_h, actual.V_h, allowedDeltaDegrees)
	assert.InDelta(excpected.Axis, actual.Axis, allowedDeltaRadians)
	assert.InDelta(excpected.Exc, actual.Exc, allowedDeltaRadians)
	assert.InDelta(DegreesFromRadians(excpected.Nu), DegreesFromRadians(actual.Nu), allowedDeltaDegrees)
}

var date4, _ = ParseDateJSON("1972-01-25T06:27")
var movement4 = NewMovement(Input{
	Tau1:  -12.55,
	Tau2:  -0.39,
	V_avg: Average([]float64{56.36, 60.91, 55.40}),
	Date:  date4,
})

func TestOrbit4(t *testing.T) {
	CheckOrbit(t, &Movement[float64]{
		Lambda_apex:  RadiansFromDegrees(214.00),
		A:            RadiansFromDegrees(21.17),
		Z_avg:        RadiansFromDegrees(42.30),
		Delta:        RadiansFromDegrees(8.81),
		Alpha:        RadiansFromDegrees(197.62),
		Beta:         RadiansFromDegrees(15.04),
		Lambda:       RadiansFromDegrees(192.77),
		Lambda_theta: RadiansFromDegrees(171.23),
		Beta_deriv:   RadiansFromDegrees(26.75),
		Inc:          RadiansFromDegrees(145.37),
		Wmega:        RadiansFromDegrees(290.42),
		Omega:        RadiansFromDegrees(304.36),
		V_g:          57.94,
		V_h:          33.41,
		Axis:         1.29,
		Exc:          0.64,
		Nu:           RadiansFromDegrees(249.59),
	}, movement4)
}

var date8, _ = ParseDateJSON("1972-01-25T07:07")
var movement8 = NewMovement(Input{
	Tau1:  -17.29,
	Tau2:  -9.95,
	V_avg: Average([]float64{32.61, 33.58, 32.22}),
	Date:  date8,
})

func TestOrbit8(t *testing.T) {
	CheckOrbit(t, &Movement[float64]{
		Lambda_apex:  RadiansFromDegrees(214.03),
		A:            RadiansFromDegrees(37.42),
		Z_avg:        RadiansFromDegrees(32.53),
		Delta:        RadiansFromDegrees(21.18),
		Alpha:        RadiansFromDegrees(201.12),
		Beta:         RadiansFromDegrees(27.72),
		Lambda:       RadiansFromDegrees(190.71),
		Lambda_theta: RadiansFromDegrees(107.45),
		Beta_deriv:   RadiansFromDegrees(51.83),
		Inc:          RadiansFromDegrees(77.10),
		Wmega:        RadiansFromDegrees(344.47),
		Omega:        RadiansFromDegrees(304.39),
		V_g:          32.98,
		V_h:          19.52,
		Axis:         0.62,
		Exc:          0.75,
		Nu:           RadiansFromDegrees(195.53),
	}, movement8)
}

// var date9, _ = ParseDate("1972-01-25T07:21")
// var movement9 = NewMovement(Input{
// 	Tau1:  -9.67,
// 	Tau2:  3.50,
// 	V_avg: Average([]float64{64.16, 67.17, 48.93}),
// 	Date:  date9,
// })

// func TestOrbit9(t *testing.T) {
// 	CheckOrbit(t, &Movement[float64]{
// 		Lambda_apex:  RadiansFromDegrees(214.04),
// 		A:            RadiansFromDegrees(10.48),
// 		Z_avg:        RadiansFromDegrees(39.83),
// 		Delta:        RadiansFromDegrees(9.99),
// 		Alpha:        RadiansFromDegrees(218.64),
// 		Beta:         RadiansFromDegrees(23.82),
// 		Lambda:       RadiansFromDegrees(212.77),
// 		Lambda_theta: RadiansFromDegrees(211.55),
// 		Beta_deriv:   RadiansFromDegrees(40.84),
// 		Inc:          RadiansFromDegrees(139.12),
// 		Wmega:        RadiansFromDegrees(184.46),
// 		Omega:        RadiansFromDegrees(304.39),
// 		V_g:          67.54,
// 		V_h:          41.70,
// 		Axis:         13.99,
// 		Exc:          0.93,
// 		Nu:           RadiansFromDegrees(355.54),
// 	}, movement9)
// }

var date11, _ = ParseDateJSON("1972-01-25T07:47")
var movement11 = NewMovement(Input{
	Tau1:  -12.92,
	Tau2:  -67.86,
	V_avg: Average([]float64{28.31, 27.21, 28.37}),
	Date:  date11,
})

func TestOrbit11(t *testing.T) {
	CheckOrbit(t, &Movement[float64]{
		Lambda_apex:  RadiansFromDegrees(214.06),
		A:            RadiansFromDegrees(102.60),
		Z_avg:        RadiansFromDegrees(68.56),
		Delta:        RadiansFromDegrees(23.93),
		Alpha:        RadiansFromDegrees(147.41),
		Beta:         RadiansFromDegrees(10.15),
		Lambda:       RadiansFromDegrees(141.47),
		Lambda_theta: RadiansFromDegrees(84.07),
		Beta_deriv:   RadiansFromDegrees(8.18),
		Inc:          RadiansFromDegrees(12.52),
		Wmega:        RadiansFromDegrees(302.98),
		Omega:        RadiansFromDegrees(304.41),
		V_g:          27.95,
		V_h:          34.62,
		Axis:         1.47,
		Exc:          0.78,
		Nu:           RadiansFromDegrees(237.03),
	}, movement11)
}

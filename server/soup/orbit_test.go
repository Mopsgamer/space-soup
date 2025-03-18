package soup

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CheckOrbit(assert *assert.Assertions, input *Input, expected *Movement) {
	var actual, err = NewMovement(*input)
	if !assert.NoError(err) {
		return
	}

	valueOfExpected := reflect.ValueOf(expected)
	valueOfActual := reflect.ValueOf(actual)
	InDelta := func(delta float64, propName string) bool {
		e := reflect.Indirect(valueOfExpected).FieldByName(propName).Float()
		a := reflect.Indirect(valueOfActual).FieldByName(propName).Float()
		return assert.InDelta(e, a, delta, fmt.Sprintf(
			"Checking \"%s\". Expected %.4f (%.2f°), got %.4f (%.2f°).",
			propName, e, DegreesFromRadians(e), a, DegreesFromRadians(a),
		))
	}
	InDelta(allowedDeltaRadians, "Lambda_apex")
	InDelta(allowedDeltaRadians, "A")
	InDelta(allowedDeltaRadians, "Z_avg")
	InDelta(allowedDeltaRadians, "Delta")
	InDelta(allowedDeltaRadians, "Alpha")
	InDelta(allowedDeltaRadians, "Beta")
	InDelta(allowedDeltaRadians, "Lambda")
	InDelta(allowedDeltaRadians, "Lambda_theta")
	InDelta(allowedDeltaRadians, "Beta_deriv")
	InDelta(allowedDeltaRadians, "Inc")
	InDelta(allowedDeltaRadians, "Wmega")
	InDelta(allowedDeltaRadians, "Omega")
	InDelta(allowedDeltaSpeed, "V_g")
	InDelta(allowedDeltaSpeed, "V_h")
	InDelta(allowedDeltaRadians, "Axis")
	InDelta(allowedDeltaRadians, "Exc")
	InDelta(allowedDeltaRadians, "Nu")
}

func TestOrbit3(t *testing.T) {
	assert := assert.New(t)
	var date, err = ParseDateJSON("1972-01-25T06:07")
	if !assert.NoError(err) {
		return
	}
	CheckOrbit(
		assert,
		&Input{
			Tau1:  -12.7572,
			Tau2:  -17.5536,
			V_avg: Average([]float64{33.858, 33.832, 33.965}),
			Date:  date,
		},
		&Movement{
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
		},
	)
}

func TestOrbit4(t *testing.T) {
	assert := assert.New(t)
	var date, err = ParseDateJSON("1972-01-25T06:27")
	if !assert.NoError(err) {
		return
	}
	CheckOrbit(
		assert,
		&Input{
			Tau1:  -12.55,
			Tau2:  -0.39,
			V_avg: Average([]float64{56.36, 60.91, 55.40}),
			Date:  date,
		},
		&Movement{
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
		},
	)
}

func TestOrbit8(t *testing.T) {
	assert := assert.New(t)
	var date, err = ParseDateJSON("1972-01-25T07:07")
	if !assert.NoError(err) {
		return
	}
	CheckOrbit(
		assert,
		&Input{
			Tau1:  -17.29,
			Tau2:  -9.95,
			V_avg: Average([]float64{32.61, 33.58, 32.22}),
			Date:  date,
		},
		&Movement{
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
		},
	)
}

func TestOrbit9(t *testing.T) {
	assert := assert.New(t)
	var date, err = ParseDateJSON("1972-01-25T07:21")
	if !assert.NoError(err) {
		return
	}
	CheckOrbit(
		assert,
		&Input{
			Tau1:  -9.67,
			Tau2:  3.50,
			V_avg: Average([]float64{64.16, 67.17, 48.93}),
			Date:  date,
		},
		&Movement{
			Lambda_apex:  RadiansFromDegrees(214.04),
			A:            RadiansFromDegrees(10.48),
			Z_avg:        RadiansFromDegrees(39.83),
			Delta:        RadiansFromDegrees(9.99),
			Alpha:        RadiansFromDegrees(218.64),
			Beta:         RadiansFromDegrees(23.82),
			Lambda:       RadiansFromDegrees(212.77),
			Lambda_theta: RadiansFromDegrees(211.55),
			Beta_deriv:   RadiansFromDegrees(40.84),
			Inc:          RadiansFromDegrees(139.12),
			Wmega:        RadiansFromDegrees(184.46),
			Omega:        RadiansFromDegrees(304.39),
			V_g:          67.54,
			V_h:          41.70,
			Axis:         13.99,
			Exc:          0.93,
			Nu:           RadiansFromDegrees(355.54),
		},
	)
}

func TestOrbit11(t *testing.T) {
	assert := assert.New(t)
	var date, err = ParseDateJSON("1972-01-25T07:47")
	if !assert.NoError(err) {
		return
	}
	CheckOrbit(
		assert,
		&Input{
			Tau1:  -12.92,
			Tau2:  -67.86,
			V_avg: Average([]float64{28.31, 27.21, 28.37}),
			Date:  date,
		},
		&Movement{
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
		},
	)
}

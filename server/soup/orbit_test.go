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
			"Checking \"%s\". Expected %.4f (%.4f°), got %.4f (%.4f°).",
			propName, e, DegreesFromRadians(e), a, DegreesFromRadians(a),
		))
	}
	InDelta(allowedDeltaRadians, "Lambda_apex")
	InDelta(allowedDeltaRadians, "H")
	InDelta(allowedDeltaRadians, "A")
	InDelta(allowedDeltaRadians, "Z_avg")
	InDelta(allowedDeltaRadians, "Delta")
	InDelta(allowedDeltaRadians, "Alpha")
	InDelta(allowedDeltaRadians, "Beta")
	InDelta(allowedDeltaRadians, "Lambda")
	InDelta(allowedDeltaRadians, "Lambda_deriv")
	InDelta(allowedDeltaRadians, "Beta_deriv")
	InDelta(allowedDeltaRadians, "Inc")
	InDelta(allowedDeltaRadians, "Wmega")
	InDelta(allowedDeltaRadians, "Omega")
	InDelta(allowedDeltaSpeed, "V_g")
	InDelta(allowedDeltaSpeed, "V_h")
	InDelta(allowedDeltaAxis, "Axis")
	InDelta(allowedDeltaExc, "Exc")
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
			Dist:  200,
			Tau1:  -12.7572,
			Tau2:  -17.5536,
			V_avg: Average([]float64{33.858, 33.832, 33.965}),
			Date:  date,
		},
		&Movement{
			Lambda_apex:  RadiansFromDegrees(213.991),
			H:            RadiansFromDegrees(78.15),
			A:            RadiansFromDegrees(59.827),
			Z_avg:        RadiansFromDegrees(27.004),
			Delta:        RadiansFromDegrees(31.756),
			Alpha:        RadiansFromDegrees(179.072),
			Beta:         RadiansFromDegrees(28.514),
			Lambda:       RadiansFromDegrees(165.364),
			Lambda_deriv: RadiansFromDegrees(99.318),
			Beta_deriv:   RadiansFromDegrees(33.340),
			Inc:          RadiansFromDegrees(57.259),
			Wmega:        RadiansFromDegrees(320.130),
			Omega:        RadiansFromDegrees(304.344),
			V_g:          34.236,
			V_h:          29.737,
			Axis:         0.9662,
			Exc:          0.7571,
			Nu:           RadiansFromDegrees(219.870),
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
			Dist:  261,
			Tau1:  -12.5536,
			Tau2:  -0.3927,
			V_avg: Average([]float64{56.3600, 60.9080, 55.3980}),
			Date:  date,
		},
		&Movement{
			Lambda_apex:  RadiansFromDegrees(214.004),
			H:            RadiansFromDegrees(81.17),
			A:            RadiansFromDegrees(21.169),
			Z_avg:        RadiansFromDegrees(42.302),
			Delta:        RadiansFromDegrees(8.810),
			Alpha:        RadiansFromDegrees(197.620),
			Beta:         RadiansFromDegrees(15.042),
			Lambda:       RadiansFromDegrees(192.772),
			Lambda_deriv: RadiansFromDegrees(171.229),
			Beta_deriv:   RadiansFromDegrees(26.764),
			Inc:          RadiansFromDegrees(145.374),
			Wmega:        RadiansFromDegrees(290.415),
			Omega:        RadiansFromDegrees(304.357),
			V_g:          57.938,
			V_h:          33.412,
			Axis:         1.2931,
			Exc:          0.6391,
			Nu:           RadiansFromDegrees(249.585),
		},
	)
}

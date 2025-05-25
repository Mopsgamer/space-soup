package soup

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CheckOrbit(assert *assert.Assertions, input Input, expected Movement) {
	actual := NewMovement(input)
	if !assert.NoError(actual.Fail) {
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
	InDelta(allowedDeltaA, "A")
	InDelta(allowedDeltaZ_avg, "Z_avg")
	InDelta(allowedDeltaDelta, "Delta")
	InDelta(allowedDeltaAlpha, "Alpha")
	InDelta(allowedDeltaBeta, "Beta")
	InDelta(allowedDeltaLambda, "Lambda")
	InDelta(allowedDeltaLambdaDer, "Lambda_deriv")
	InDelta(allowedDeltaBetaDer, "Beta_deriv")
	InDelta(allowedDeltaInc, "Inc")
	InDelta(allowedDeltaWmega, "Wmega")
	InDelta(allowedDeltaOmega, "Omega")
	InDelta(allowedDeltaVg, "V_g")
	InDelta(allowedDeltaVh, "V_h")
	InDelta(allowedDeltaAxis, "Axis")
	InDelta(allowedDeltaExc, "Exc")
	InDelta(allowedDeltaNu, "Nu")
}

func TestOrbit800016(t *testing.T) {
	assert := assert.New(t)
	date, _ := ParseDateJSON("1978-01-02T10:28")
	CheckOrbit(
		assert,
		Input{
			Id:    800016,
			Tau1:  -12.9214,
			Tau2:  -13.5505,
			V_avg: Average([]float64{49.037, 49.604}),
			Date:  date,
		},
		Movement{
			A:            RadiansFromDegrees(51.248),
			Lambda_apex:  RadiansFromDegrees(193.337),
			Z_avg:        RadiansFromDegrees(38.540),
			Delta:        RadiansFromDegrees(19.833),
			Alpha:        RadiansFromDegrees(219.011),
			Beta:         RadiansFromDegrees(33.152),
			Lambda:       RadiansFromDegrees(209.186),
			Lambda_deriv: RadiansFromDegrees(240.753),
			Beta_deriv:   RadiansFromDegrees(60.410),
			Inc:          RadiansFromDegrees(111.022),
			Wmega:        RadiansFromDegrees(87.743),
			Omega:        RadiansFromDegrees(283.346),
			V_g:          50.882,
			V_h:          31.998,
			Axis:         1.1365,
			Exc:          0.3846,
			Nu:           RadiansFromDegrees(92.257),
		},
	)
}

func BenchmarkNewMovement(b *testing.B) {
	var date, err = ParseDateJSON("1972-01-25T06:07")
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for range b.N {
		NewMovement(Input{
			Tau1:  -12.7572,
			Tau2:  -17.5536,
			V_avg: Average([]float64{33.858, 33.832, 33.965}),
			Date:  date,
		})
	}
}

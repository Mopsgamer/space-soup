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
	InDelta(allowedDeltaRadians, "Lambda_apex")
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

package soup

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type MovementTest struct {
	Input    Input
	Expected *Movement
	Actual   *Movement
	// 0 - Success, 1 - Acceptable, 2 - Not acceptable
	AssertionResult *MovementAssertion
}

func CheckOrbitList() (result []MovementTest, err error) {
	result = []MovementTest{}
	bytes, err := os.ReadFile("ORB_72.txt")
	if err != nil {
		return
	}

	linesOut := strings.Split(string(bytes), "\n")

	bytes, err = os.ReadFile("orb-72.txt")
	if err != nil {
		return
	}

	linesIn := strings.Split(string(bytes), "\n")[1:]
	inputList := map[int]Input{}
	for _, line := range linesIn {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		dt := strings.Split(fields[2], ".")
		date := time.Date(1900+Int(dt[2]), time.Month(Int(dt[0])), Int(dt[1]), Int(fields[4]), Int(fields[5]), 0, 0, time.UTC)

		speedList := []float64{}
		for _, v := range []float64{Float64(fields[6]), Float64(fields[7]), Float64(fields[8])} {
			if v > 999. {
				continue
			}
			speedList = append(speedList, v)
		}

		id := Int(fields[0])
		input := Input{
			Id:    &id,
			Dist:  Float64(fields[1]),
			Tau1:  Float64(fields[9]),
			Tau2:  Float64(fields[10]),
			Date:  date,
			V_avg: Average(speedList),
		}
		inputList[Int(fields[0])] = input
	}

	for _, line := range linesOut {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n := Int(fields[0]) - 200000
		input, ok := inputList[n]
		if !ok {
			continue
		}
		actual, err := NewMovement(input)
		if err != nil && *input.Id < 50 {
			fmt.Printf("%d-th got error: %s\n", 200000+*input.Id, err)
		}

		entry := MovementTest{
			Input:           input,
			Actual:          &Movement{},
			Expected:        &Movement{},
			AssertionResult: &MovementAssertion{},
		}

		entry.Actual.Lambda_apex = DegreesFromRadians(actual.Lambda_apex)
		entry.Actual.H = actual.H
		entry.Actual.A = DegreesFromRadians(actual.A)
		entry.Actual.Z_avg = DegreesFromRadians(actual.Z_avg)
		entry.Actual.Delta = DegreesFromRadians(actual.Delta)
		entry.Actual.Alpha = DegreesFromRadians(actual.Alpha)
		entry.Actual.Beta = DegreesFromRadians(actual.Beta)
		entry.Actual.Lambda = DegreesFromRadians(actual.Lambda)
		entry.Actual.Lambda_deriv = DegreesFromRadians(actual.Lambda_deriv)
		entry.Actual.Beta_deriv = DegreesFromRadians(actual.Beta_deriv)
		entry.Actual.Inc = DegreesFromRadians(actual.Inc)
		entry.Actual.Wmega = DegreesFromRadians(actual.Wmega)
		entry.Actual.Omega = DegreesFromRadians(actual.Omega)
		entry.Actual.V_g = actual.V_g
		entry.Actual.V_h = actual.V_h
		entry.Actual.Axis = actual.Axis
		entry.Actual.Exc = actual.Exc
		entry.Actual.Nu = DegreesFromRadians(actual.Nu)

		entry.Expected.Lambda_apex = Float64(fields[5])
		entry.Expected.H = Float64(fields[9])
		entry.Expected.A = Float64(fields[10])
		entry.Expected.Z_avg = Float64(fields[11])
		entry.Expected.Delta = Float64(fields[12])
		entry.Expected.Alpha = Float64(fields[13])
		entry.Expected.Beta = Float64(fields[14])
		entry.Expected.Lambda = Float64(fields[15])
		entry.Expected.Lambda_deriv = Float64(fields[16])
		entry.Expected.Beta_deriv = Float64(fields[17])
		entry.Expected.Inc = Float64(fields[18])
		entry.Expected.Wmega = Float64(fields[19])
		entry.Expected.Omega = Float64(fields[20])
		entry.Expected.V_g = Float64(fields[21])
		entry.Expected.V_h = Float64(fields[22])
		entry.Expected.Axis = Float64(fields[23])
		entry.Expected.Exc = Float64(fields[24])
		entry.Expected.Nu = Float64(fields[25])

		valueOfExpected := reflect.ValueOf(entry.Expected)
		valueOfActual := reflect.ValueOf(entry.Actual)
		InDelta := func(delta float64, propName string) uint {
			e := reflect.Indirect(valueOfExpected).FieldByName(propName).Float()
			a := reflect.Indirect(valueOfActual).FieldByName(propName).Float()
			return InDelta(e, a, delta)
		}

		entry.AssertionResult.Lambda_apex = InDelta(allowedDeltaDegrees, "Lambda_apex")
		entry.AssertionResult.H = InDelta(allowedDeltaDegrees, "H")
		entry.AssertionResult.A = InDelta(allowedDeltaDegrees, "A")
		entry.AssertionResult.Z_avg = InDelta(allowedDeltaDegrees, "Z_avg")
		entry.AssertionResult.Delta = InDelta(allowedDeltaDegrees, "Delta")
		entry.AssertionResult.Alpha = InDelta(allowedDeltaDegrees, "Alpha")
		entry.AssertionResult.Beta = InDelta(allowedDeltaDegrees, "Beta")
		entry.AssertionResult.Lambda = InDelta(allowedDeltaDegrees, "Lambda")
		entry.AssertionResult.Lambda_deriv = InDelta(allowedDeltaDegrees, "Lambda_deriv")
		entry.AssertionResult.Beta_deriv = InDelta(allowedDeltaDegrees, "Beta_deriv")
		entry.AssertionResult.Inc = InDelta(allowedDeltaDegrees, "Inc")
		entry.AssertionResult.Wmega = InDelta(allowedDeltaDegrees, "Wmega")
		entry.AssertionResult.Omega = InDelta(allowedDeltaDegrees, "Omega")
		entry.AssertionResult.V_g = InDelta(allowedDeltaSpeed, "V_g")
		entry.AssertionResult.V_h = InDelta(allowedDeltaSpeed, "V_h")
		entry.AssertionResult.Axis = InDelta(allowedDeltaAxis, "Axis")
		entry.AssertionResult.Exc = InDelta(allowedDeltaExc, "Exc")
		entry.AssertionResult.Nu = InDelta(allowedDeltaDegrees, "Nu")
		result = append(result, entry)
	}
	return
}

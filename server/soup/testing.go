package soup

import (
	"os"
	"strings"
	"time"
)

type TestEntry[T any] struct {
	Expected T
	Actual   T
	// 0 - Success, 1 - Acceptable, 2 - Not acceptable
	AssertionResult uint
}

func CheckOrbitList() (result []Movement[TestEntry[float64]], err error) {
	result = []Movement[TestEntry[float64]]{}
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

		input := Input{
			Dist:  Int(fields[1]),
			Tau1:  Float64(fields[9]),
			Tau2:  Float64(fields[10]),
			Date:  date,
			V_avg: Average([]float64{Float64(fields[6]), Float64(fields[7]), Float64(fields[8])}),
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
		if err != nil {
			continue
		}
		entry := Movement[TestEntry[float64]]{
			Lambda_apex:  TestEntry[float64]{Actual: actual.Lambda_apex, Expected: Float64(fields[5])},
			A:            TestEntry[float64]{Actual: actual.A, Expected: Float64(fields[10])},
			Z_avg:        TestEntry[float64]{Actual: actual.Z_avg, Expected: Float64(fields[11])},
			Delta:        TestEntry[float64]{Actual: actual.Delta, Expected: Float64(fields[12])},
			Alpha:        TestEntry[float64]{Actual: actual.Alpha, Expected: Float64(fields[13])},
			Beta:         TestEntry[float64]{Actual: actual.Beta, Expected: Float64(fields[14])},
			Lambda:       TestEntry[float64]{Actual: actual.Lambda, Expected: Float64(fields[15])},
			Lambda_theta: TestEntry[float64]{Actual: actual.Lambda_theta, Expected: Float64(fields[16])},
			Beta_deriv:   TestEntry[float64]{Actual: actual.Beta_deriv, Expected: Float64(fields[17])},
			Inc:          TestEntry[float64]{Actual: actual.Inc, Expected: Float64(fields[18])},
			Wmega:        TestEntry[float64]{Actual: actual.Wmega, Expected: Float64(fields[19])},
			Omega:        TestEntry[float64]{Actual: actual.Omega, Expected: Float64(fields[20])},
			V_g:          TestEntry[float64]{Actual: actual.V_g, Expected: Float64(fields[21])},
			V_h:          TestEntry[float64]{Actual: actual.V_h, Expected: Float64(fields[22])},
			Axis:         TestEntry[float64]{Actual: actual.Axis, Expected: Float64(fields[23])},
			Exc:          TestEntry[float64]{Actual: actual.Exc, Expected: Float64(fields[24])},
			Nu:           TestEntry[float64]{Actual: actual.Nu, Expected: Float64(fields[25])},
		}
		result = append(result, entry)
	}
	return
}

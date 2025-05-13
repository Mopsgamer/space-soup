package soup

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "embed"

	"github.com/gofiber/fiber/v3/log"
)

var idInc = 800000

//go:embed ORB_78.PAR
var fileExpected string

//go:embed ORB_1978-save.csv
var fileInput string

type MovementTest struct {
	Input           Input
	Expected        Movement
	Actual          Movement
	AssertionResult MovementAssertion
}

func CheckOrbitList() (tests []MovementTest, err error) {
	start := time.Now()
	fnStart := start
	var sinceStart, sincefnStart, sincefnStart2 time.Duration
	stop := func() {
		sinceStart = time.Since(start)
		sincefnStart2 += sinceStart
	}

	linesOut := strings.Split(fileExpected, "\n")
	stop()
	log.Infof("Read answers file and split %d lines: %v", len(linesOut), sinceStart)

	start = time.Now()
	reader := csv.NewReader(strings.NewReader(fileInput))
	rowsIn, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse fileInput: %w", err)
	}
	rowsIn = rowsIn[1:]
	stop()
	log.Infof("Read input file and split %d lines: %v", len(rowsIn), sinceStart)

	start = time.Now()
	inputList := map[int]Input{}
	for _, fields := range rowsIn {
		if len(fields) == 0 {
			continue
		}
		dt := strings.Split(fields[2], "/")
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
			Id:    id,
			Tau1:  Float64(fields[9]),
			Tau2:  Float64(fields[10]),
			Date:  date,
			V_avg: Average(speedList),
		}
		inputList[Int(fields[0])] = input
	}
	stop()
	log.Infof("Parse input %d lines: %v", len(rowsIn), sinceStart)

	start = time.Now()
	fieldsOut := [][]string{}
	validInputList := []Input{}
	for _, line := range linesOut {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		fieldsOut = append(fieldsOut, fields)
		n := Int(fields[0]) - idInc
		input, ok := inputList[n]
		if !ok {
			continue
		}
		validInputList = append(validInputList, input)
	}
	stop()
	log.Infof("Parse answers %d lines: %v", len(linesOut), sinceStart)

	start = time.Now()
	actualList := []Movement{}
	for _, input := range validInputList {
		actual := NewMovement(input)
		if actual.Fail != nil {
			fmt.Printf("%d-th got error: %s\n", idInc+input.Id, actual.Fail)
		}
		actualList = append(actualList, actual)
	}
	stop()
	log.Infof("Calculate %d orbits: %v (%v/1)", len(validInputList), sinceStart, time.Duration(float64(sinceStart)/float64(len(validInputList))))

	start = time.Now()
	tests = make([]MovementTest, len(actualList))
	for n, actual := range actualList {
		input := validInputList[n]
		fields := fieldsOut[n]
		entry := MovementTest{
			Input:           input,
			Actual:          actual,
			Expected:        Movement{},
			AssertionResult: MovementAssertion{},
		}

		entry.Expected.Lambda_apex = RadiansFromDegrees(Float64(fields[5]))
		// entry.Expected.H = RadiansFromDegrees(Float64(fields[9]))
		entry.Expected.A = RadiansFromDegrees(Float64(fields[10]))
		entry.Expected.Z_avg = RadiansFromDegrees(Float64(fields[11]))
		entry.Expected.Delta = RadiansFromDegrees(Float64(fields[12]))
		entry.Expected.Alpha = RadiansFromDegrees(Float64(fields[13]))
		entry.Expected.Beta = RadiansFromDegrees(Float64(fields[14]))
		entry.Expected.Lambda = RadiansFromDegrees(Float64(fields[15]))
		entry.Expected.Lambda_deriv = RadiansFromDegrees(Float64(fields[16]))
		entry.Expected.Beta_deriv = RadiansFromDegrees(Float64(fields[17]))
		entry.Expected.Inc = RadiansFromDegrees(Float64(fields[18]))
		entry.Expected.Wmega = RadiansFromDegrees(Float64(fields[19]))
		entry.Expected.Omega = RadiansFromDegrees(Float64(fields[20]))
		entry.Expected.V_g = Float64(fields[21])
		entry.Expected.V_h = Float64(fields[22])
		entry.Expected.Axis = Float64(fields[23])
		entry.Expected.Exc = Float64(fields[24])
		entry.Expected.Nu = RadiansFromDegrees(Float64(fields[25]))

		valueOfExpected := reflect.ValueOf(entry.Expected)
		valueOfActual := reflect.ValueOf(entry.Actual)
		InDelta := func(delta float64, propName string) TestResult {
			e := reflect.Indirect(valueOfExpected).FieldByName(propName).Float()
			a := reflect.Indirect(valueOfActual).FieldByName(propName).Float()
			return InDelta(e, a, delta)
		}

		entry.AssertionResult.Lambda_apex = InDelta(allowedDeltaRadians, "Lambda_apex")
		// entry.AssertionResult.H = InDelta(allowedDeltaRadians, "H")
		entry.AssertionResult.A = InDelta(allowedDeltaRadians, "A")
		entry.AssertionResult.Z_avg = InDelta(allowedDeltaRadians, "Z_avg")
		entry.AssertionResult.Delta = InDelta(allowedDeltaRadians, "Delta")
		entry.AssertionResult.Alpha = InDelta(allowedDeltaRadians, "Alpha")
		entry.AssertionResult.Beta = InDelta(allowedDeltaRadians, "Beta")
		entry.AssertionResult.Lambda = InDelta(allowedDeltaRadians, "Lambda")
		entry.AssertionResult.Lambda_deriv = InDelta(allowedDeltaRadians, "Lambda_deriv")
		entry.AssertionResult.Beta_deriv = InDelta(allowedDeltaRadians, "Beta_deriv")
		entry.AssertionResult.Inc = InDelta(allowedDeltaRadians, "Inc")
		entry.AssertionResult.Wmega = InDelta(allowedDeltaRadians, "Wmega")
		entry.AssertionResult.Omega = InDelta(allowedDeltaRadians, "Omega")
		entry.AssertionResult.V_g = InDelta(allowedDeltaSpeed, "V_g")
		entry.AssertionResult.V_h = InDelta(allowedDeltaSpeed, "V_h")
		entry.AssertionResult.Axis = InDelta(allowedDeltaAxis, "Axis")
		entry.AssertionResult.Exc = InDelta(allowedDeltaExc, "Exc")
		entry.AssertionResult.Nu = InDelta(allowedDeltaRadians, "Nu")
		tests[n] = entry
	}
	stop()
	sincefnStart = time.Since(fnStart)
	log.Infof("Test %d orbits: %v (%v/1)", len(actualList), sinceStart, time.Duration(float64(sinceStart)/float64(len(validInputList))))

	log.Infof("Summary: %v (%v/print all)", sincefnStart, (sincefnStart - sincefnStart2).Abs())
	return
}

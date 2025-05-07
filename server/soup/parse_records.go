package soup

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidRowFormat = errors.New("invalid row format")
)

var notnum = regexp.MustCompile(`[^\s0-9,.]`)

func ParseRecords(records [][]string) ([]MovementTest, error) {
	var movementTestList []MovementTest
	for i, record := range records {
		if len(record) < 4 {
			return nil, ErrInvalidRowFormat
		}
	ForLine:
		for i, field := range record {
			if i == 2 {
				continue
			}
			if notnum.MatchString(field) {
				continue ForLine
			}
		}
		error := new(error)
		*error = nil
		tau1 := Float64Err(record[0], error)
		tau2 := Float64Err(record[1], error)
		date, err := ParseDateJSON(record[2])
		if err != nil {
			return nil, errors.Join(*error, err)
		}
		vList := []float64{}
		for _, vString := range record[3:] {
			v := Float64Err(vString, error)
			vList = append(vList, v)
		}
		if *error != nil {
			return nil, *error
		}

		input := Input{
			Id:    i + 1,
			Tau1:  tau1,
			Tau2:  tau2,
			V_avg: Average(vList),
			Date:  date,
		}
		movement := NewMovement(input)
		movementTest := MovementTest{
			Input:  input,
			Actual: movement,
		}
		if movement.Fail == nil {
			movementTestList = append(movementTestList, movementTest)
		}
	}
	return movementTestList, nil
}

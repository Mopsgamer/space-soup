package soup

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
)

var (
	ErrInvalidRowFormat    = errors.New("invalid row format")
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

var notnum = regexp.MustCompile(`[^\s0-9,.]`)

func NewMovementTestsFromFile(formFile *multipart.FileHeader, filetype FileType) ([]MovementTest, error) {
	var movementTestList []MovementTest

	filetype = filetype.DecideFileName(formFile.Filename)

	file, err := formFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rows [][]string

	switch string(filetype) {
	case "csv":
		fallthrough
	case "tsv":
		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1
		reader.TrimLeadingSpace = true
		if filetype[0] == 't' { // filetype == tsv
			reader.Comma = '\t'
		}
		rows, err = reader.ReadAll()
	case "xlsx":
		// var xlFile *excelize.File
		// xlFile, err = excelize.OpenReader(file)
		// if err != nil {
		// 	return nil, err
		// }
		// sheetName := xlFile.GetSheetName(1)
		// rows, err = xlFile.GetRows(sheetName)
		fallthrough
	default:
		return nil, errors.Join(ErrUnsupportedFileType, fmt.Errorf("file type is '%s'", filetype))
	}

	if err != nil {
		return nil, err
	}
	movementTestList, err = ParseRecords(rows)
	if err != nil {
		return nil, err
	}

	return movementTestList, nil
}

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

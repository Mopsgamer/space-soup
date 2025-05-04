package model_http

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/Mopsgamer/space-soup/server/soup"
	"github.com/xuri/excelize/v2"
)

var (
	ErrInvalidRowFormat = errors.New("invalid row format")
)

type OrbitInputFile struct {
	FileType FileType `form:"file-type"`
}

func (p *OrbitInputFile) MovementTestList(pFile multipart.FileHeader) ([]soup.MovementTest, error) {
	var movementTestList []soup.MovementTest

	p.FileType = p.FileType.DecideFileName(pFile.Filename)

	file, err := pFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	switch string(p.FileType) {
	case "csv":
		fallthrough
	case "tsv":
		reader := csv.NewReader(file)
		if p.FileType[0] == 't' { // filetype == tsv
			reader.Comma = '\t'
		}
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		movementTestList, err = parseRecords(records)
		if err != nil {
			return nil, err
		}
	case "xlsx":
		xlFile, err := excelize.OpenReader(file)
		if err != nil {
			return nil, err
		}
		sheetName := xlFile.GetSheetName(1)
		rows, err := xlFile.GetRows(sheetName)
		if err != nil {
			return nil, err
		}
		movementTestList, err = parseRecords(rows)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Join(ErrUnsupportedFileType, fmt.Errorf("file type is '%s'", p.FileType))
	}

	return movementTestList, nil
}

func parseRecords(records [][]string) ([]soup.MovementTest, error) {
	var movementTestList []soup.MovementTest
	for _, record := range records {
		if len(record) < 3 {
			return nil, ErrInvalidRowFormat
		}
		error := new(error)
		*error = nil
		tau1 := soup.Float64Err(record[0], error)
		tau2 := soup.Float64Err(record[1], error)
		date, err := soup.ParseDateJSON(record[2])
		if err != nil {
			return nil, errors.Join(*error, err)
		}
		vList := []float64{}
		for _, vString := range record[3:] {
			v := soup.Float64Err(vString, error)
			vList = append(vList, v)
		}
		if *error != nil {
			return nil, *error
		}

		input := soup.Input{
			Tau1:  tau1,
			Tau2:  tau2,
			V_avg: soup.Average(vList),
			Date:  date,
		}
		movement := soup.NewMovement(input)
		movementTest := soup.MovementTest{
			Input:  input,
			Actual: movement,
		}
		if movement.Fail == nil {
			movementTestList = append(movementTestList, movementTest)
		}
	}
	return movementTestList, nil
}

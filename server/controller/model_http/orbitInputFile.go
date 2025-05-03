package model_http

import (
	"encoding/csv"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/Mopsgamer/space-soup/server/soup"
	"github.com/gofiber/fiber/v3/log"
	"github.com/xuri/excelize/v2"
)

var (
	ErrInvalidRowFormat = errors.New("invalid row format")
)

type OrbitInputFile struct {
	FileType FileType `form:"file-type"`
}

func (p *OrbitInputFile) MovementList(pFile multipart.FileHeader) ([]soup.Movement, error) {
	var movements []soup.Movement

	if p.FileType.ShouldDecide() {
		ext := FileType(strings.ToLower(filepath.Ext(pFile.Filename)))
		switch ext {
		case ".csv":
			p.FileType = "csv"
		case ".tsv":
			p.FileType = "tsv"
		case ".xlsx":
			p.FileType = "xlsx"
		default:
			log.Info("can not decide for ext", ext)
		}
	}

	file, err := pFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	switch string(p.FileType) {
	case "csv":
	case "tsv":
		reader := csv.NewReader(file)
		if p.FileType[0] == 't' { // filetype == tsv
			reader.Comma = '\t'
		}
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		movements, err = parseRecords(records)
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
		movements, err = parseRecords(rows)
		if err != nil {
			return nil, err
		}
	default:
		log.Info("ubknown file type", p.FileType)
		return nil, ErrUnsupportedFileType
	}

	return movements, nil
}

func parseRecords(records [][]string) ([]soup.Movement, error) {
	var movements []soup.Movement
	for _, record := range records {
		if len(record) < 3 {
			return nil, ErrInvalidRowFormat
		}
		error := new(error)
		*error = nil
		tau1 := soup.Float64Err(record[0], error)
		tau2 := soup.Float64Err(record[1], error)
		vList := []float64{}
		for _, vString := range record[4:] {
			v := soup.Float64Err(vString, error)
			vList = append(vList, v)
		}
		if *error != nil {
			return nil, *error
		}
		date, err := soup.ParseDateJSON(record[3])
		if err != nil {
			return nil, errors.Join(*error, err)
		}

		input := soup.Input{
			Tau1:  tau1,
			Tau2:  tau2,
			V_avg: soup.Average(vList),
			Date:  date,
		}
		movement := soup.NewMovement(input)
		if movement.Fail == nil {
			movements = append(movements, movement)
		}
	}
	return movements, nil
}

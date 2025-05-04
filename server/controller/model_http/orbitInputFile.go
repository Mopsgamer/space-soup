package model_http

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/Mopsgamer/space-soup/server/soup"
	"github.com/xuri/excelize/v2"
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
		movementTestList, err = soup.ParseRecords(records)
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
		movementTestList, err = soup.ParseRecords(rows)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Join(ErrUnsupportedFileType, fmt.Errorf("file type is '%s'", p.FileType))
	}

	return movementTestList, nil
}

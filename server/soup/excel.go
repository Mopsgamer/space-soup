package soup

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

func NewFileExcelBytes(movementList []MovementTest) ([]byte, error) {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(0)

	bordered, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "Times New Roman"},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	headers := []string{"#", "V_avg", "Tau1", "Tau2", "Lambda_apex", "A", "Z_avg", "Delta", "Alpha", "Beta", "Lambda", "Lambda_deriv", "Beta_deriv", "Inc", "Wmega", "Omega", "V_g", "V_h", "Axis", "Exc", "Nu"}
	for colIndex, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIndex+1, 1)
		file.SetCellValue(sheetName, cell, header)
		file.SetCellStyle(sheetName, cell, cell, bordered)
	}

	colIndex := 0
	nextColCell := func(row int) string {
		colIndex += 1
		cell, _ := excelize.CoordinatesToCellName(colIndex, row)
		return cell
	}

	for rowIndex, movement := range movementList {
		row := rowIndex + 2 // Start from row 2 as row 1 is the header
		colIndex = 0
		setCellInt := func(v int) {
			coord := nextColCell(row)
			file.SetCellInt(sheetName, coord, v)
			file.SetCellStyle(sheetName, coord, coord, bordered)
		}
		setCellFloat := func(v float64) {
			coord := nextColCell(row)
			file.SetCellFloat(sheetName, coord, v, 2, 32)
			file.SetCellStyle(sheetName, coord, coord, bordered)
		}
		setCellInt(movement.Input.Id)
		setCellFloat(movement.Input.V_avg)
		setCellFloat(movement.Input.Tau1)
		setCellFloat(movement.Input.Tau2)
		setCellFloat(movement.Actual.Lambda_apex)
		setCellFloat(movement.Actual.A)
		setCellFloat(movement.Actual.Z_avg)
		setCellFloat(movement.Actual.Delta)
		setCellFloat(movement.Actual.Alpha)
		setCellFloat(movement.Actual.Beta)
		setCellFloat(movement.Actual.Lambda)
		setCellFloat(movement.Actual.Lambda_deriv)
		setCellFloat(movement.Actual.Beta_deriv)
		setCellFloat(movement.Actual.Inc)
		setCellFloat(movement.Actual.Wmega)
		setCellFloat(movement.Actual.Omega)
		setCellFloat(movement.Actual.V_g)
		setCellFloat(movement.Actual.V_h)
		setCellFloat(movement.Actual.Axis)
		setCellFloat(movement.Actual.Exc)
		setCellFloat(movement.Actual.Nu)
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return nil, err
	}

	excelBytes := buf.Bytes()
	return excelBytes, nil
}

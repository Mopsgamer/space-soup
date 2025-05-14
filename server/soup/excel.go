package soup

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

func NewFileExcelBytes(movementList []MovementTest) ([]byte, error) {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(0)

	bordered, _ := file.NewStyle(&excelize.Style{Border: []excelize.Border{
		{Type: "left", Style: 1, Color: "000000"},
		{Type: "top", Style: 1, Color: "000000"},
		{Type: "right", Style: 1, Color: "000000"},
		{Type: "bottom", Style: 1, Color: "000000"},
	}})
	// Add header row
	headers := []string{"V_avg", "Tau1", "Tau2", "Lambda_apex", "A", "Z_avg", "Delta", "Alpha", "Beta", "Lambda", "Lambda_deriv", "Beta_deriv", "Inc", "Wmega", "Omega", "V_g", "V_h", "Axis", "Exc", "Nu"}
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
		setCell := func(v float64) {
			coord := nextColCell(row)
			file.SetCellFloat(sheetName, coord, v, 2, 32)
			file.SetCellStyle(sheetName, coord, coord, bordered)
		}
		// V_avg	Tau1	Tau2	Lambda_apex	A	Z_avg	Delta	Alpha	Beta	Lambda	Lambda_deriv	Beta_deriv	Inc	Wmega	Omega	V_g	V_h	Axis	Exc	Nu
		setCell(movement.Input.V_avg)
		setCell(movement.Input.Tau1)
		setCell(movement.Input.Tau2)
		setCell(movement.Actual.Lambda_apex)
		setCell(movement.Actual.A)
		setCell(movement.Actual.Z_avg)
		setCell(movement.Actual.Delta)
		setCell(movement.Actual.Alpha)
		setCell(movement.Actual.Beta)
		setCell(movement.Actual.Lambda)
		setCell(movement.Actual.Lambda_deriv)
		setCell(movement.Actual.Beta_deriv)
		setCell(movement.Actual.Inc)
		setCell(movement.Actual.Wmega)
		setCell(movement.Actual.Omega)
		setCell(movement.Actual.V_g)
		setCell(movement.Actual.V_h)
		setCell(movement.Actual.Axis)
		setCell(movement.Actual.Exc)
		setCell(movement.Actual.Nu)
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return nil, err
	}

	excelBytes := buf.Bytes()
	return excelBytes, nil
}

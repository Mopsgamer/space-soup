package soup

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

func NewFileExcelBytes(movementList []MovementTest) ([]byte, error) {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(0)

	colIndex := 0
	nextColCell := func(row int) string {
		colIndex += 1
		cell, _ := excelize.CoordinatesToCellName(colIndex, row)
		return cell
	}
	for rowIndex, movement := range movementList {
		colIndex = 0
		// V_avg	Tau1	Tau2	Lambda_apex	A	Z_avg	Delta	Alpha	Beta	Lambda	Lambda_deriv	Beta_deriv	Inc	Wmega	Omega	V_g	V_h	Axis	Exc	Nu
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Input.V_avg, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Input.Tau1, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Input.Tau2, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Lambda_apex, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.A, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Z_avg, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Delta, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Alpha, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Beta, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Lambda, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Lambda_deriv, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Beta_deriv, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Inc, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Wmega, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Omega, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.V_g, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.V_h, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Axis, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Exc, 2, 32)
		file.SetCellFloat(sheetName, nextColCell(rowIndex), movement.Actual.Nu, 2, 32)
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return nil, err
	}

	excelBytes := buf.Bytes()
	return excelBytes, nil
}

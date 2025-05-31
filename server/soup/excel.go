package soup

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/xuri/excelize/v2"
)

func NewFileExcelBytes(movementList []MovementTest, noDelta bool) ([]byte, error) {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(0)

	normalCfg := &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Font: &excelize.Font{Family: "Times New Roman", Color: "000000", Size: 14},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	}
	normal, _ := file.NewStyle(normalCfg)

	headerCfg := &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Font: &excelize.Font{
			Family: "Times New Roman",
			Color:  "000000",
			Size:   14,
			Bold:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	}
	header, _ := file.NewStyle(headerCfg)

	headers := []string{
		"№",
		"№ метеороїда",
		"Середня швидкість (V_avg)",
		"Часова затримка 1 (Tau1)",
		"Часова затримка 2 (Tau2)",
		"Довгота апексу (Lambda_apex)",
		"Азимут (A)",
		"Зенітний кут радіанта (Z_avg)",
		"Схиляння радіанта (Delta)",
		"Пряме сходження радіанта (Alpha)",
		"Екліптична широта радіанта (Beta)",
		"Екліптична довгота радіанта (Lambda)",
		"Довгота істинного радіанта (Lambda_deriv)",
		"Широта істинного радіанта (Beta_deriv)",
		"Нахил орбіти (Inc)",
		"Довгота висхідного вузла (Wmega)",
		"Аргумент перигелію (Omega)",
		"Геоцентрична швидкість (V_g)",
		"Геліоцентрична швидкість (V_h)",
		"Велика піввісь (Axis)",
		"Ексцентриситет (Exc)",
		"Істинна аномалія (Nu)",
	}
	for colIndex, label := range headers {
		colIndex += 1
		cell, _ := excelize.CoordinatesToCellName(colIndex, 1)
		file.SetCellValue(sheetName, cell, label)
		file.SetCellStyle(sheetName, cell, cell, header)

		cell2, _ := excelize.CoordinatesToCellName(colIndex, 2)
		file.SetCellValue(sheetName, cell2, colIndex)
		file.SetCellStyle(sheetName, cell2, cell2, header)
	}

	colIndex := 0
	nextColCell := func(row int) string {
		colIndex += 1
		cell, _ := excelize.CoordinatesToCellName(colIndex, row)
		return cell
	}

	for rowIndex, movement := range movementList {
		rowIndex += 1
		row := rowIndex + 2
		colIndex = 0
		setCellInt := func(v int) {
			coord := nextColCell(row)
			file.SetCellInt(sheetName, coord, v)
			file.SetCellStyle(sheetName, coord, coord, normal)
		}
		setCellFloat := func(v float64, suffix string) {
			coord := nextColCell(row)
			str := strings.Replace(fmt.Sprintf("%.2f%s", v, suffix), ".", ",", 1)
			file.SetCellValue(sheetName, coord, str)
			file.SetCellStyle(sheetName, coord, coord, normal)
		}
		setCellFloatDelta := func(actual, expected float64, suffix string) {
			if noDelta {
				setCellFloat(actual, suffix)
				return
			}
			coord := nextColCell(row)
			delta := math.Abs(actual - expected)
			str := strings.Replace(fmt.Sprintf("%.2f%s ± %.2f%s", actual, suffix, delta, suffix), ".", ",", -1)
			file.SetCellValue(sheetName, coord, str)
			file.SetCellStyle(sheetName, coord, coord, normal)
		}
		setCellInt(rowIndex + 1)
		setCellInt(movement.Input.Id)
		setCellFloat(movement.Input.V_avg, "")
		setCellFloat(movement.Input.Tau1, "")
		setCellFloat(movement.Input.Tau2, "")
		setCellFloatDelta(movement.Actual.Lambda_apex, movement.Expected.Lambda_apex, "°")
		setCellFloatDelta(movement.Actual.A, movement.Expected.A, "°")
		setCellFloatDelta(movement.Actual.Z_avg, movement.Expected.Z_avg, "°")
		setCellFloatDelta(movement.Actual.Delta, movement.Expected.Delta, "°")
		setCellFloatDelta(movement.Actual.Alpha, movement.Expected.Alpha, "°")
		setCellFloatDelta(movement.Actual.Beta, movement.Expected.Beta, "°")
		setCellFloatDelta(movement.Actual.Lambda, movement.Expected.Lambda, "°")
		setCellFloatDelta(movement.Actual.Lambda_deriv, movement.Expected.Lambda_deriv, "°")
		setCellFloatDelta(movement.Actual.Beta_deriv, movement.Expected.Beta_deriv, "°")
		setCellFloatDelta(movement.Actual.Inc, movement.Expected.Inc, "°")
		setCellFloatDelta(movement.Actual.Wmega, movement.Expected.Wmega, "°")
		setCellFloatDelta(movement.Actual.Omega, movement.Expected.Omega, "°")
		setCellFloatDelta(movement.Actual.V_g, movement.Expected.V_g, "")
		setCellFloatDelta(movement.Actual.V_h, movement.Expected.V_h, "")
		setCellFloatDelta(movement.Actual.Axis, movement.Expected.Axis, "")
		setCellFloatDelta(movement.Actual.Exc, movement.Expected.Exc, "")
		setCellFloatDelta(movement.Actual.Nu, movement.Expected.Nu, "°")
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return nil, err
	}

	excelBytes := buf.Bytes()
	return excelBytes, nil
}

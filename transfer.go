package dynamic

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Transfer[T any] struct {
}

func (t *Transfer[T]) values(file *excelize.File, sheet string) ([]CellValue, error) {
	var values []CellValue
	rows, err := file.Rows(sheet)
	if err != nil {
		return values, err
	}

	curY := 0
	for rows.Next() {
		curY++

		cols, err := rows.Columns()
		if err != nil {
			continue
		}

		size := len(cols)
		for i := 0; i < size; i++ {
			value := CellValue{
				X:     i + 1,
				Y:     curY,
				Value: cols[i],
			}

			values = append(values, value)
		}
	}
	rows.Close()

	return values, nil
}

func (t *Transfer[T]) Read(tree *Tree[T], file *excelize.File, sheet string) error {
	_, err := file.GetSheetIndex(sheet)
	if err != nil {
		return err
	}

	values, err := t.values(file, sheet)
	if err != nil {
		return err
	}

	for idx, value := range values {
		fmt.Printf("%3d - (%s, %s)\n", idx, fmt.Sprintf("%s%d", numberToLetters(value.X), value.Y), value.Value)
	}

	return nil
}

package dynamic

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Sheet[T any] struct {
	tree         *Tree[T]
	dataRowStart int
	mapdValues   map[string]*CellValue
}

func (c *Sheet[T]) values(file *excelize.File, sheet string) (map[string]*CellValue, error) {
	if c.mapdValues != nil {
		return c.mapdValues, nil
	}

	c.mapdValues = make(map[string]*CellValue)

	rows, err := file.Rows(sheet)
	if err != nil {
		return nil, err
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

			key := fmt.Sprintf("%s%d", numberToLetters(value.X), value.Y)
			c.mapdValues[key] = &value
		}
	}
	rows.Close()
	c.build_relation()

	return c.mapdValues, nil
}

func (c *Sheet[T]) build_relation() {
	for _, cellValue := range c.mapdValues {
		nextCell := cellValue.NextCell()
		parentCell := cellValue.ParentCell()

		if nextCellValue, ok := c.mapdValues[nextCell]; ok {
			cellValue.Next = nextCellValue
			nextCellValue.Prev = cellValue
		}

		if parentCellValue, ok := c.mapdValues[parentCell]; ok {
			cellValue.Parent = parentCellValue
		}
	}
}

func (c *Sheet[T]) Read(tree *Tree[T], file *excelize.File, sheet string) error {
	_, err := file.GetSheetIndex(sheet)
	if err != nil {
		return err
	}
	c.tree = tree

	values, err := c.values(file, sheet)
	if err != nil {
		return err
	}

	for idx, value := range values {
		fmt.Printf("%s - %s\n", idx, value.Value)
	}

	return nil
}

// DataRowStart
// 数据行起始，是根据传入数据结构的判断
func (c *Sheet[T]) DataRowStart() int {
	return c.tree.MaxLevel() + 1
}

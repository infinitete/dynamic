package dynamic

import (
	"fmt"
	"log"
	"slices"

	"strings"

	"github.com/xuri/excelize/v2"
)

type Sheet[T any] struct {
	file         *excelize.File
	sheet        string
	tree         *Tree[T]
	dataRowStart int
	mapdValues   map[string]*CellValue

	// aliasCells
	// 合并关系,key是被合并对象,value是主索引
	// 例如合并: [A1, A3], 那么产生两个合并关系:
	// A2:A1和A3:A1
	aliasCells map[string]string
}

func (c *Sheet[T]) buildAlias() {
	if c.aliasCells != nil {
		return
	}
	mergedCells, err := c.file.GetMergeCells(c.sheet)
	if err != nil {
		return
	}

	c.aliasCells = make(map[string]string)

	for _, cell := range mergedCells {
		start := cell.GetStartAxis()
		end := cell.GetEndAxis()
		cells := fillCells(start, end)
		alias := cells[0]
		for i := 1; i < len(cells); i++ {
			c.aliasCells[cells[i]] = alias
		}
	}
}

func (c *Sheet[T]) values() (map[string]*CellValue, error) {
	if c.mapdValues != nil {
		return c.mapdValues, nil
	}

	c.mapdValues = make(map[string]*CellValue)

	rows, err := c.file.Rows(c.sheet)
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

	for key, aliasKey := range c.aliasCells {
		cell, ok := c.mapdValues[key]
		if !ok {
			// 自动填充
			cell = &CellValue{}
			c.mapdValues[key] = cell
		}
		alias, ok := c.mapdValues[aliasKey]
		if !ok {
			continue
		}
		cell.Alias = alias
	}

	rows.Close()
	c.build_relation()

	return c.mapdValues, nil
}

// build_relation
// 构建非空关系树
func (c *Sheet[T]) build_relation() {
	for _, cellValue := range c.mapdValues {
		nextCell := cellValue.NextCell()
		parentCell := cellValue.ParentCell()

		if cellValue.Next == nil {
			if nextCellValue, ok := c.mapdValues[nextCell]; ok {
				cellValue.Next = nextCellValue
				nextCellValue.Prev = cellValue
			}
		}

		// 注意
		// 此处是向坐上查找
		if parentCellValue, ok := c.mapdValues[parentCell]; ok {
			if parentCellValue.Alias != nil {
				cellValue.Parent = parentCellValue.Alias
			} else {
				cellValue.Parent = parentCellValue
			}
		}
	}
}

func (c *Sheet[T]) Read(tree *Tree[T], file *excelize.File, sheet string) error {
	c.file = file
	c.sheet = sheet

	_, err := file.GetSheetIndex(sheet)
	if err != nil {
		return err
	}
	c.tree = tree
	c.buildAlias()

	cellValues, err := c.values()
	if err != nil {
		return err
	}

	log.Printf("N3: %#v", c.mapdValues["N3"])

	for key, value := range cellValues {
		if value.Y > c.tree.MaxLevel() || value.Value == "" {
			continue
		}
		paths := value.Paths()
		slices.Reverse(paths)
		fmt.Printf("[%s]: %s\n", key, strings.Join(paths, "->"))
	}

	return nil
}

// DataRowStart
// 数据行起始，是根据传入数据结构的判断
func (c *Sheet[T]) DataRowStart() int {
	return c.tree.MaxLevel() + 1
}

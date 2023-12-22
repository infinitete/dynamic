package dynamic

import (
	"fmt"
	"reflect"

	"github.com/xuri/excelize/v2"
)

func NewRenderer[T any](file *excelize.File, sheet string) (*Renderer[T], error) {
	var t T
	if reflect.TypeOf(t).Kind() != reflect.Struct {
		return nil, fmt.Errorf("typeof %T is not sutrct", t)
	}

	return &Renderer[T]{
		file:  file,
		sheet: sheet,
	}, nil
}

type Renderer[T any] struct {
	file   *excelize.File
	sheet  string
	parser *Parser[T]
	tree   *Tree[T]
}

func (r *Renderer[T]) Render(data []T) {
	if r.parser == nil {
		r.parser = &Parser[T]{}

		tree, _ := r.parser.Parse()
		r.tree = tree
		r.tree.ParseValues(data)
	}

	r.file.NewSheet(r.sheet)
	headStyleId, _ := r.file.NewStyle(&headerStyle)
	bodyStyleId, _ := r.file.NewStyle(&bodyStyle)

	for _, meta := range r.tree.metas {

		start := fmt.Sprintf("%s%d", numberToLetters(meta.StartX), meta.StartY)
		end := fmt.Sprintf("%s%d", numberToLetters(meta.EndX), meta.EndY)
		if start != end {
			r.file.MergeCell(r.sheet, start, end)
		}
		r.file.SetCellStr(r.sheet, start, meta.Node.Title)
		r.file.SetCellStyle(r.sheet, start, end, headStyleId)

		if len(meta.rows) == 0 {
			continue
		}

		meta.Render(func(coor string, value TypedValue) error {
			if err := r.file.SetCellValue(r.sheet, coor, value.Value); err != nil {
				return err
			}
			return r.file.SetCellStyle(r.sheet, coor, coor, bodyStyleId)
		})
	}
}

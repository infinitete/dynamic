package dynamic

import (
	"log"

	"github.com/xuri/excelize/v2"
)

type Reader[T any] struct {
	parser *Parser[T]
	sheet  *Sheet[T]
}

func NewReader[T any]() *Reader[T] {
	parser := Parser[T]{}
	transfer := Sheet[T]{}

	return &Reader[T]{
		parser: &parser,
		sheet:  &transfer,
	}
}

func (r *Reader[T]) Read(file *excelize.File, sheet string) []T {
	idx, err := file.GetSheetIndex(sheet)
	if idx == -1 {
		return nil
	}

	r.parser = &Parser[T]{}
	tree, err := r.parser.Parse()
	if err != nil {
		return nil
	}

	err = r.sheet.Read(tree, file, sheet)
	if err != nil {
		return nil
	}

	log.Printf("数据起始行：%d", r.sheet.DataRowStart())

	return nil
}

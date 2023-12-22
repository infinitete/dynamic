package dynamic

import (
	"github.com/xuri/excelize/v2"
)

type Reader[T any] struct {
	parser   *Parser[T]
	transfer *Transfer[T]
}

func NewReader[T any]() *Reader[T] {
	parser := Parser[T]{}
	transfer := Transfer[T]{}

	return &Reader[T]{
		parser:   &parser,
		transfer: &transfer,
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

	err = r.transfer.Read(tree, file, sheet)
	if err != nil {
		return nil
	}

	return nil
}

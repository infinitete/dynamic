package dynamic

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Reader
// 实现从excel中读取数据到指定结构体，并返回结构体数组
type Reader[T any] struct {

	// parser
	// 结构体解析器
	parser *Parser[T]

	// sheet
	// excel读取器
	sheet *Sheet[T]
}

func NewReader[T any]() (*Reader[T], error) {
	parser := Parser[T]{}
	tree, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	sheet := Sheet[T]{
		headerLevel: tree.MaxLevel(),
	}

	return &Reader[T]{
		parser: &parser,
		sheet:  &sheet,
	}, nil
}

// Read
// 从Excel读取并解析数据到对应结构体
func (r *Reader[T]) Read(file *excelize.File, sheet string) []T {
	idx, err := file.GetSheetIndex(sheet)
	if idx == -1 {
		return nil
	}

	_, err = r.sheet.Read(file, sheet)
	if err != nil {
		return nil
	}

	// 向下匹配法
	// 判断结构体节点是否和excel节点匹配的步骤
	// 1. 判断当前节点的名称和excel节点的值是否一致
	// 2. 如果有子节点，那么逐个判断子节点的值和excel对应子节点的值是否一致
	// 3. 如果子节点存在子节点，那么重复2-3步骤，直至完成匹配
	// 4. 如果有节点不匹配，那么跳出匹配
	// 5. 完成匹配
	headers := r.sheet.GetHeaderCellValues()
	for _, header := range headers {
		if header.Alias != nil {
			continue
		}

		nodes := r.parser.tree.FindNodesByTag(header.Y, header.Value)
		if len(nodes) == 0 || !r.fullMatch(nodes[0], header) {
			continue
		}
		nodes[0].offsetX = header.X
	}

	var nodeValues = map[*Node]map[int]string{}

	for _, meta := range r.parser.tree.metas {
		if len(meta.Node.Children) > 0 {
			continue
		}
		nodeValues[meta.Node] = make(map[int]string)
	}

	var startY = r.parser.tree.maxLevel
	var depth = r.sheet.depth
	for node, values := range nodeValues {
		for i := startY + 1; i <= depth; i++ {
			cell := fmt.Sprintf("%s%d", numberToLetters(node.offsetX), i)
			cellValue, ok := r.sheet.mapdValues[cell]
			if !ok {
				continue
			}
			values[i] = cellValue.Value
		}
	}
	var values = make([]T, depth-startY)
	for node, value := range nodeValues {
		for i := startY + 1; i <= depth; i++ {
			var t = values[i-startY-1]
			err := r.setStructValue(&t, node.Kind, r.parser.tree.paths(node), value[i])
			if err != nil {
				log.Printf("Error: %s", err.Error())
			}
			values[i-startY-1] = t
		}
	}

	return values
}

// match
// 判断一个从结构体中解析的node和从excel中读取的元素是否一致
// 一致的条件:
// 1. 层级一致
// 2. node的标签(title)和元素的值是否相等
func (r *Reader[T]) match(node *Node, cell *CellValue) bool {
	if node == nil || cell == nil {
		return false
	}

	if node.Level != cell.Y {
		return false
	}

	return node.Title == cell.Value
}

// fullMatch
// 判断节点和excel的节点是否全匹配
// 全匹配的条件：
// 1. 当前节点匹配,层级、值一致、子节点数量一致
// 2. 所有子节点匹配
func (r *Reader[T]) fullMatch(node *Node, cell *CellValue) bool {
	if node == nil || cell == nil {
		return false
	}

	if node.Level != cell.Y {
		return false
	}

	if node.Title != cell.Value {
		return false
	}

	// 子节点完全匹配
	children := r.sheet.GetChildrenCellValues(cell)
	if len(children) != len(node.Children) {
		return false
	}

	for _, child := range children {
		nodeChild := node.FindChildByTitle(child.Value)
		if nodeChild == nil {
			return false
		}
		if !r.fullMatch(nodeChild, child) {
			return false
		}
	}

	return true
}

// TODO
func (r *Reader[T]) setStructValue(t *T, kind reflect.Kind, paths []string, value string) error {
	var valueOf reflect.Value
	var typeOf reflect.Type
	valueOf = reflect.ValueOf(t).Elem()
	typeOf = reflect.TypeOf(t).Elem()

	for i := 0; i < len(paths); i++ {
		for {
			if typeOf.Kind() == reflect.Ptr {
				typeOf = typeOf.Elem()
				continue
			}
			break
		}

		_, ok := typeOf.FieldByName(paths[i])
		if !ok {
			return fmt.Errorf("A - path [%s] is not valid, full path is: [%s], value is: [%s]", strings.Join(paths[0:i+1], "->"), strings.Join(paths, "->"), value)
		}

		if valueOf.Kind() == reflect.Ptr && valueOf.IsNil() {
			valueOf = reflect.New(typeOf).Elem()
		}

		valueOf = valueOf.FieldByName(paths[i])
		if !valueOf.IsValid() {
			return fmt.Errorf("B - path [%s] is not valid, full path is: [%s], value is: [%s]", strings.Join(paths[0:i+1], "->"), strings.Join(paths, "->"), value)
		}
		typeOf = valueOf.Type()
	}

	if valueOf.Kind() != kind {
		return fmt.Errorf("path [%s] is no valid: expected [%s], [%s] read from struct [%T]", strings.Join(paths, "->"), kind, valueOf.Kind(), t)
	}

	switch kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			err = fmt.Errorf("value of path [%s] no set: %s", strings.Join(paths, "->"), err.Error())
		}
		valueOf.SetInt(intValue)
	case reflect.Float32,
		reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			err = fmt.Errorf("value of path [%s] no set: %s", strings.Join(paths, "->"), err.Error())
		}
		valueOf.SetFloat(floatValue)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		intValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			err = fmt.Errorf("value of path [%s] no set: %s", strings.Join(paths, "->"), err.Error())
		}
		valueOf.SetUint(intValue)
	default:
		valueOf.SetString(value)
	}

	log.Printf("Set Value: [%s] -> %v", strings.Join(paths, "->"), valueOf)

	return nil
}

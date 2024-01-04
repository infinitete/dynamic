package dynamic

import (
	"fmt"
	"reflect"
)

type TypedValue struct {
	Paths []string
	Kind  reflect.Kind
	Value any
}

type CellValue struct {
	X     int
	Y     int
	Value string

	Alias  *CellValue // 合并对象，添加合并对象的目的是为了根据结构构建结构树
	Prev   *CellValue // 上一个元素,如果有合并对象，那么他的上一个元素是合并对象的上一个元素
	Next   *CellValue // 下一个元素,如果有合并对象，那么他的下一个元素是合并对象的下一个元素
	Parent *CellValue // 父级元素,如果有合并对象，那么他的父级元素是合并对象的父级元素
}

func (c CellValue) Cell() string {
	return fmt.Sprintf("%s%d", numberToLetters(c.X), c.Y)
}

func (c CellValue) PrevCell() string {
	if c.X == 1 {
		return ""
	}

	return fmt.Sprintf("%s%d", numberToLetters(c.X-1), c.Y)
}

func (c CellValue) NextCell() string {
	return fmt.Sprintf("%s%d", numberToLetters(c.X+1), c.Y)
}

func (c CellValue) ParentCell() string {
	if c.Y == 1 {
		return ""
	}

	return fmt.Sprintf("%s%d", numberToLetters(c.X), c.Y-1)
}

func (c CellValue) Paths() []string {
	var paths = []string{c.Value}
	var parent = c.Parent
	for {
		if parent == nil {
			break
		}
		paths = append(paths, parent.Value)
		parent = parent.Parent
	}
	return paths
}

// Meta
// 可渲染元数据
type Meta struct {

	// rows
	// 数据
	rows []TypedValue

	// 节点基本信息
	Node *Node

	// 数据类型
	Kind reflect.Kind

	// Paths
	// 渲染路径
	// 渲染路径包含了所有父节点的Field
	// 例如一个结构体如下：
	/*
	   type A struct {
	      Field string `xlsx:"col:字段"`
	   }
	   type B struct {
	      BField string `xlsx:"col:另一个字段`
	      AField A `xlsx:"col:原来的字段`
	   }

	   // 实际渲染的内容
	   type C struct {
	      CField int `xlsx:"col:C字段"`
	      BField B `xlxs:"col:B字段"`
	   }
	*/
	// 展开C结构体如下：
	/*
		C:
		 CField:
		   BField:
			     BField
		     AField:
			       A
	*/
	// 那么A的路径就是[CField, BField, AField, A]
	Paths []string

	// StartX
	// 列所在坐标
	StartX int

	// StartY
	// 起始行所在坐标
	StartY int

	// EndX
	// 合并列坐标
	EndX int

	// 合并行坐标
	EndY int

	// currentY
	// 当前渲染行
	currentY int
}

type RenderFun func(coor string, value TypedValue) error

func (m *Meta) Render(fn RenderFun) {
	dataSize := len(m.rows)
	if dataSize == 0 {
		return
	}

	for cur := 0; cur < dataSize; cur++ {
		coor := fmt.Sprintf("%s%d", numberToLetters(m.StartX), m.EndY+cur+1)
		fn(coor, m.rows[cur])
	}
}

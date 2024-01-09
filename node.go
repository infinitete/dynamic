package dynamic

import (
	"fmt"
	"reflect"
	"strings"
)

const ignore = "-"

type Node struct {
	parent   *Node
	index    int
	offsetX  int
	Field    string
	Title    string
	Level    int
	Depth    *int
	Kind     reflect.Kind
	Children []*Node
}

func (node Node) Y() int {
	return node.Level
}

// Rows
// 获取节点所占行数
func (node Node) Rows() int {
	return *node.Depth - node.Level + 1
}

// Cols
// 获取该节点所占列数
func (node Node) Cols() int {
	if len(node.Children) == 0 {
		return 1
	}
	cols := 0
	for _, child := range node.Children {
		cols = cols + child.Cols()
	}
	return cols
}

func (node Node) CanMergeRows() bool {
	return node.Rows() > 1
}

func (node Node) ToCellValues() CellValue {
	value := CellValue{
		X:     node.offsetX,
		Y:     node.Level,
		Value: node.Title,
	}

	return value
}

func (node Node) FindChildByTitle(title string) *Node {
	for _, child := range node.Children {
		if child.Title == title {
			return child
		}
	}

	return nil
}

// Paths
// 获取节点的值路径
func (node Node) Paths() []string {
	var paths = []string{node.Title}
	n := node.parent
	for {
		if n == nil {
			break
		}
		paths = append(paths, n.Title)
		n = n.parent
	}
	return paths
}

type Parser[T any] struct {
	tree *Tree[T]
}

func (p *Parser[T]) Tree() *Tree[T] {
	tree, _ := p.Parse()
	return tree
}

// Parse
// 解析指定的类型，并返回一颗解析树
// 类型不是结构体，那么返回错误
// 如果类型中包含指针，那么解析返回错误
// 为什么不解析指针：
// 这个库的目标是读取和写入excel，excel中任何一个节点的值都有确定类型和值
func (p *Parser[T]) Parse() (*Tree[T], error) {
	if p.tree != nil {
		return p.tree, nil
	}

	var t T
	depth := 1

	nodes, err := p.parseType(reflect.TypeOf(t), nil, &depth, 1)
	if err != nil {
		return nil, err
	}

	tree := Tree[T]{Nodes: nodes}
	_ = tree.Metas()
	p.tree = &tree

	return p.tree, nil
}

func (p *Parser[T]) parseType(typeOf reflect.Type, parent *Node, depth *int, level int) ([]*Node, error) {
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	fields := typeOf.NumField()
	nodes := []*Node{}

	*depth = level

	cur := 0
	for i := 0; i < fields; i++ {
		field := typeOf.Field(i)

		if field.Type.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("type %s mus not be contains  pointer field", typeOf.Name())
		}

		tag := field.Tag.Get("xlsx")
		if tag == ignore {
			continue
		}
		node := &Node{
			index:  cur,
			parent: parent,
			Level:  level,
			Depth:  depth,
		}
		cur++
		node.Field = field.Name
		cases := strings.Split(tag, ",")
		for _, c := range cases {
			kv := strings.Split(c, ":")
			if kv[0] == "col" {
				if len(kv) == 2 {
					node.Title = kv[1]
				}
				break // get col only
			}
		}
		if node.Title == "" {
			node.Title = node.Field
		}

		node.Kind = field.Type.Kind()

		if node.Kind == reflect.Struct {
			children, err := p.parseType(field.Type, node, depth, level+1)
			if err != nil {
				return nil, err
			}
			node.Children = children
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

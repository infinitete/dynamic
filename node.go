package dynamic

import (
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

func (node Node) Rows() int {
	return *node.Depth - node.Level + 1
}

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

type Parser[T any] struct {
	tree *Tree[T]
}

func (p Parser[T]) Tree() *Tree[T] {
	tree, _ := p.Parse()
	return tree
}

func (p Parser[T]) Parse() (*Tree[T], error) {
	if p.tree != nil {
		return p.tree, nil
	}

	var t T
	depth := 1

	tree := Tree[T]{Nodes: p.parseType(reflect.TypeOf(t), nil, &depth, 1)}
	_ = tree.Metas()
	p.tree = &tree

	return p.tree, nil
}

func (p Parser[T]) parseType(typeOf reflect.Type, parent *Node, depth *int, level int) []*Node {
	fields := typeOf.NumField()
	nodes := []*Node{}

	*depth = level

	cur := 0
	for i := 0; i < fields; i++ {
		field := typeOf.Field(i)
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
			node.Children = p.parseType(field.Type, node, depth, level+1)
		}
		nodes = append(nodes, node)
	}

	return nodes
}

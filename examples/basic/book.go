package main

import (
	"fmt"
	"strings"

	"github.com/infinitete/dynamic"
)

type Names struct {
	English string `xlsx:"col:英文名"`
	Chinese string `xlsx:"col:中文名"`
	Franch  string `xlsx:"col:法文名"`
}

type Pointer struct {
	X int `xlsx:"col:x"`
	Y int `xlsx:"col:y"`
}

type Author struct {
	FirstName string `xlsx:"col:姓"`
	LastName  string `xlsx:"col:名"`
}

type Title struct {
	MainTitle Names `xlsx:"col:主标题"`
	SubTitle  Names `xlsx:"col:副标题"`
}

type Bookmark struct {
	Index int     `xlsx:"col:序号"`
	Page  int     `xlsx:"col:页数"`
	Start Pointer `xlsx:"col:起始"`
	End   Pointer `xlsx:"col:终点"`
}

type Book struct {
	Title    Title    `xlsx:"col:书名"`
	Author   Author   `xlsx:"col:作者"`
	Bookmark Bookmark `xlsx:"col:书签"`
	Remark   string   `xlsx:"col:备注"`
}

func main() {
	parser := dynamic.Parser{}
	tree := parser.Parse(Book{})

	for _, node := range tree.Nodes {
		printNode(tree, node)
	}

	for _, meta := range tree.Metas() {
		fmt.Printf("%s\n", strings.Join(meta.Paths, "->"))
	}
}

func printNode(tree dynamic.Tree, node *dynamic.Node) {
	prefix := make([]string, node.Level)
	fmt.Printf("%s%s(%d:%d)\n", strings.Join(prefix, "    "), node.Title, tree.OffsetX(node), node.Cols())
	for _, n := range node.Children {
		printNode(tree, n)
	}
}

package dynamic

import (
	"fmt"
	"strings"
)

func NewRenderer[T any](title string) Renderer[T] {
	return Renderer[T]{title: title}
}

type Renderer[T any] struct {
	title  string
	parser *Parser
	tree   *Tree
}

func (r *Renderer[T]) Render([]T) {
	if r.parser == nil {
		var t T
		r.parser = &Parser{}

		tree := r.parser.Parse(t)
		r.tree = &tree
	}

	metas := r.tree.Metas()
	for _, meta := range metas {
		r.renderMeta(meta)
	}
}

func (r *Renderer[T]) renderMeta(meta *Meta) {
	coorStart := fmt.Sprintf("%s%d", numberToLetters(meta.StartX), meta.StartY)
	coorEnd := fmt.Sprintf("%s%d", numberToLetters(meta.EndX), meta.EndY)

	fmt.Printf("%s%s(%s, %s)\n", strings.Join(make([]string, meta.Node.Level), "  "), meta.Node.Title, coorStart, coorEnd)
}

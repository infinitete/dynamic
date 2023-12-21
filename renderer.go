package dynamic

type Renderer[T any] struct {
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
}

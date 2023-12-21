package dynamic

type Tree struct {
	Nodes    []*Node
	metas    []*Meta
	maxLevel int
}

// offsetX
// 获取一个节点在一棵树中的X偏移
func (t *Tree) offsetX(node *Node) int {
	if node == nil {
		return 1
	}

	prev := t.GetPrev(node)
	if prev == nil {
		return t.offsetX(node.parent)
	}

	return t.offsetX(prev) + prev.Cols()
}

// paths
// 获取一个节点的路径
// 路径组成是一个数组，这个数组是一个按层级递增的的
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
func (t *Tree) paths(node *Node) []string {
	if node == nil {
		return []string{}
	}
	paths := []string{node.Field}
	paths = append(t.paths(t.GetParent(node)), paths...)

	return paths
}

func (t *Tree) meta(node *Node) *Meta {
	if node == nil {
		return nil
	}

	meta := &Meta{
		Node:     node,
		Paths:    t.paths(node),
		StartX:   t.OffsetX(node),
		StartY:   node.Y(),
		CurrentY: node.Y(),
	}

	meta.EndX = meta.StartX + meta.Node.Cols() - 1
	meta.EndY = meta.StartY

	return meta
}

func (t *Tree) childrenMetas(node *Node) []*Meta {
	if node == nil {
		return nil
	}

	metas := []*Meta{}
	for _, child := range node.Children {
		childMeta := t.meta(child)
		metas = append(metas, childMeta)
		if len(node.Children) > 0 {
			metas = append(metas, t.childrenMetas(child)...)
		}
	}

	return metas
}

func (t *Tree) GetParent(node *Node) *Node {
	return node.parent
}

func (t *Tree) GetPrev(node *Node) *Node {
	index := node.index
	if index == 0 {
		return nil
	}

	parent := t.GetParent(node)
	if parent == nil {
		for _, n := range t.Nodes {
			if n.index == index-1 {
				return n
			}
		}
		return nil
	}

	for _, n := range parent.Children {
		if n.index == index-1 {
			return n
		}
	}

	return nil
}

func (t *Tree) GetNext(node *Node) *Node {
	index := node.index
	if index == 0 {
		return nil
	}

	parent := t.GetParent(node)
	if parent == nil {
		size := len(t.Nodes)
		if node.index == size-1 {
			return nil
		}
		for _, n := range t.Nodes {
			if n.index == node.index+1 {
				return n
			}
		}
		return nil
	}

	size := len(parent.Children)
	if node.index == size-1 {
		return nil
	}

	for _, n := range parent.Children {
		if n.index == index+1 {
			return n
		}
	}

	return nil
}

func (t Tree) OffsetX(node *Node) int {
	if node.offsetX == 0 {
		node.offsetX = t.offsetX(node)
	}

	return node.offsetX
}

// Metas
// 获取每一个元素的可渲染元数据
func (t *Tree) Metas() []*Meta {
	if t.metas != nil {
		return t.metas
	}

	var metas []*Meta

	for _, node := range t.Nodes {
		nodeMeta := t.meta(node)
		nodeChilrenMetas := t.childrenMetas(node)
		metas = append(metas, nodeMeta)
		metas = append(metas, nodeChilrenMetas...)
	}
	t.metas = metas

	for _, meta := range metas {
		if len(meta.Node.Children) == 0 {
			meta.EndY = t.MaxLevel()
		}
	}

	return t.metas
}

func (t *Tree) MaxLevel() int {
	if t.Nodes == nil {
		return 0
	}

	if t.maxLevel == 0 {
		for _, meta := range t.metas {
			if meta.Node.Level > t.maxLevel {
				t.maxLevel = meta.Node.Level
			}
		}
	}

	return t.maxLevel
}

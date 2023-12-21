package dynamic

// Meta
// 可渲染元数据
type Meta struct {
	Node *Node

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

	// CurrentY
	// 当前渲染行
	CurrentY int
}

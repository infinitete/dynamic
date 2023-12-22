package main

import (
	"github.com/infinitete/dynamic"
	"github.com/xuri/excelize/v2"
)

type Names struct {
	Chinese string `xlsx:"col:中文名"`
	English string `xlsx:"col:英文名"`
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
	Title   string   `xlsx:"-"`
	Index   int      `xlsx:"col:序号"`
	Page    int      `xlsx:"col:页数"`
	Start   Pointer  `xlsx:"col:起始"`
	End     Pointer  `xlsx:"col:终点"`
	Readers []string `xlsx:"读者"`
}

type Book struct {
	NotPresented string   `xlsx:"-"`
	Title        Title    `xlsx:"col:书名"`
	Author       Author   `xlsx:"col:作者"`
	Bookmark     Bookmark `xlsx:"col:书签"`
	Remark       string   `xlsx:"col:备注"`
}

var book = Book{
	NotPresented: "这一列是不会渲染的",
	Title: Title{
		MainTitle: Names{
			English: "Song of Zhangsan",
			Chinese: "张三之歌",
			Franch:  "Chanson de trois",
		},
		SubTitle: Names{
			English: "Confession of an Extralegal Madman",
			Chinese: "一个法外狂徒的自白",
			Franch:  "Confession d'un maniaque extra - judiciaire",
		},
	},
	Author: Author{
		FirstName: "罗",
		LastName:  "用好",
	},
	Bookmark: Bookmark{
		Title: "虚拟书签",
		Index: 1,
		Page:  213,
		Start: Pointer{
			X: 32,
			Y: 24,
		},
		End: Pointer{
			X: 65,
			Y: 122,
		},
		Readers: []string{"张三", "李四", "王五"},
	},
	Remark: "这本书真的不错哦",
}

func books() {
	// 第一步，创建一个excel文件
	file := excelize.NewFile()

	// 第二步，申明一个渲染器渲染类型
	renderer, err := dynamic.NewRenderer[Book](file, "书本")
	if err != nil {
		return
	}

	size := 10000
	books := make([]Book, size)
	for cur := 0; cur < size; cur++ {
		book.Bookmark.Index = cur + 1
		books[cur] = book
	}

	// 第三步，渲染数据
	renderer.Render(books)

	// 第四步, 保存文件
	file.SaveAs("book.xlsx")
}

func main() {
	books()
}

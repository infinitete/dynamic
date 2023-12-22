package dynamic

import "github.com/xuri/excelize/v2"

var headerStyle = excelize.Style{
	Fill: excelize.Fill{
		Type:    "pattern",
		Color:   []string{"#e0ebf5"},
		Pattern: 1,
	},
	Font: &excelize.Font{
		Bold:      true,
		Italic:    false,
		Underline: "",
		Family:    "Times New Roman",
		Size:      12,
		Color:     "#000000",
	},
	Border: []excelize.Border{{
		Type:  "left",
		Color: "333333",
		Style: 3,
	}, {
		Type:  "top",
		Color: "333333",
		Style: 3,
	}, {
		Type:  "bottom",
		Color: "333333",
		Style: 3,
	}, {
		Type:  "right",
		Color: "333333",
		Style: 3,
	}},
	Alignment: &excelize.Alignment{
		Horizontal:      "center",
		Indent:          1,
		JustifyLastLine: true,
		ReadingOrder:    0,
		RelativeIndent:  1,
		ShrinkToFit:     true,
		Vertical:        "center",
		WrapText:        true,
	},
}

var bodyStyle = excelize.Style{
	Alignment: &excelize.Alignment{
		Horizontal:      "center",
		Indent:          1,
		JustifyLastLine: true,
		ReadingOrder:    0,
		RelativeIndent:  0,
		ShrinkToFit:     true,
		TextRotation:    0,
		Vertical:        "center",
		WrapText:        true,
	},
}

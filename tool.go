package dynamic

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func numberToLetters(num int) string {
	if num <= 0 {
		return ""
	}

	result := ""
	for num > 0 {
		// 'A'的ASCII码为65，所以减去1来得到对应的字母
		letter := string(rune('A' + ((num - 1) % 26)))
		result = letter + result
		num = (num - 1) / 26
	}

	return result
}

func lettersToNumber(letters string) int {
	base := 26
	result := 0

	letters = strings.ToUpper(letters)

	for i := range letters {
		letter := letters[i]
		letterValue := int(letter - 'A' + 1)
		result += letterValue * int(math.Pow(float64(base), float64(len(letters)-i-1)))
	}

	return result
}

// GetFieldsValue
// 通过反射获取传入参数中的所有值
func getFieldsValue(in any) map[string]TypedValue {
	valueOf := reflect.ValueOf(in)
	typeOf := reflect.TypeOf(in)

	var typedValues = []*TypedValue{}

	numField := valueOf.NumField()
	for cur := 0; cur < numField; cur++ {
		typedValues = getChildrenTypedValue(valueOf, typeOf, []string{})
	}

	values := map[string]TypedValue{}

	for _, value := range typedValues {
		key := strings.Join(value.Paths, ".")
		values[key] = *value
	}

	return values
}

func getChildrenTypedValue(valueOf reflect.Value, typeOf reflect.Type, beforePaths []string) []*TypedValue {
	values := []*TypedValue{}

	numField := valueOf.NumField()
	for cur := 0; cur < numField; cur++ {
		fieldValue := valueOf.Field(cur)
		fieldType := typeOf.Field(cur)
		if fieldType.Tag.Get("xlsx") == "-" {
			continue
		}

		paths := append(beforePaths, fieldType.Name)

		if fieldType.Type.Kind() == reflect.Struct {
			values = append(values, getChildrenTypedValue(fieldValue, fieldType.Type, paths)...)
		} else {
			values = append(values, &TypedValue{
				Paths: paths,
				Kind:  fieldType.Type.Kind(),
				Value: fieldValue.Interface(),
			})
		}
	}

	return values
}

func fillCells(startCell, endCell string) []string {
	if startCell == endCell {
		return []string{startCell}
	}

	var startXBytes = []byte{}
	var endXBytes = []byte{}
	var startYBytes = []byte{}
	var endYBytes = []byte{}

	for i := 0; i < len(startCell); i++ {
		if startCell[i] >= 'A' && startCell[i] <= 'Z' {
			startXBytes = append(startXBytes, startCell[i])
		} else {
			startYBytes = []byte(startCell)[i:]
			break
		}
	}

	for i := 0; i < len(endCell); i++ {
		if endCell[i] >= 'A' && endCell[i] <= 'Z' {
			endXBytes = append(endXBytes, endCell[i])
		} else {
			endYBytes = []byte(endCell)[i:]
			break
		}
	}

	var startX = lettersToNumber(string(startXBytes))
	var endX = lettersToNumber(string(endXBytes))
	var startY, _ = strconv.ParseInt(string(startYBytes), 10, 64)
	var endY, _ = strconv.ParseInt(string(endYBytes), 10, 64)

	var cells = make([]string, (endX-startX+1)*int((endY-startY+1)))

	var idx = 0
	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			cells[idx] = fmt.Sprintf("%s%d", numberToLetters(x), y)
			idx++
		}
	}

	return cells
}

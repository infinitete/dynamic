package dynamic

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

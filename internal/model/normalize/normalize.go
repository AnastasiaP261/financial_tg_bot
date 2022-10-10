package normalize

import (
	"strings"
	"unicode"
)

// Category делает первый символ в строке заглавным, а остальные строчными
func Category(str string) string {
	res := []rune(strings.TrimSpace(str))
	for i := range res {
		if i == 0 {
			res[i] = unicode.ToUpper(res[i])
		} else if unicode.IsLetter(res[i]) {
			res[i] = unicode.ToLower(res[i])
		}
	}

	return string(res)
}

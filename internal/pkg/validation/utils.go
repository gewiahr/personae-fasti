package validation

import (
	"slices"
	"unicode/utf8"
)

func IsValidLength(str string, min, max int) (bool, int) {
	charCount := utf8.RuneCountInString(str)
	if charCount < min {
		return false, charCount - min
	} else if charCount > max {
		return false, charCount - max
	}

	return true, 0
}

func IsValidString(str string, letters bool, numbers bool, additional []rune) bool {
	for _, r := range str {

		isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		if isLetter {
			if !letters {
				return false
			}
			continue
		}

		isNumber := (r >= '0' && r <= '9')
		if isNumber {
			if !numbers {
				return false
			}
			continue
		}

		isAdditional := slices.Contains(additional, r)
		if isAdditional {
			continue
		}

		return false

	}

	return true
}

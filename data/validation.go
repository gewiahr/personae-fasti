package data

import (
	"slices"
	"unicode/utf8"
)

func (s *Storage) ValidatePlayerUsername(username string) (bool, string) {
	if valid, count := s.isValidLength(username, 6, 30); !valid {
		if count > 0 {
			return false, "Логин не может быть больше 30 символов"
		} else if count < 0 {
			return false, "Логин не может быть меньше 6 символов"
		}
	}

	if valid := s.isValidString(username, true, true, []rune{'.', '-', '_'}); !valid {
		return false, "Логин может содержать только латинские буквы, цифры, точку, дефис или нижнее подчёркивание"
	}

	return true, "Логин валиден"
}

func (s *Storage) ValidatePlayerPassword(password string) (bool, string) {
	if valid, count := s.isValidLength(password, 8, 64); !valid {
		if count > 0 {
			return false, "Пароль не может быть больше 64 символов"
		} else if count < 0 {
			return false, "Пароль не может быть меньше 8 символов"
		}
	}

	if valid := s.isValidString(password, true, true, []rune{'.', ',', '-', '_', '!', '@', '#', '$', '%', '^', '&', '*', '+', '=', '?', '/', '\\'}); !valid {
		return false, "Пароль содержит некорректные символы - возможно введены не латинские буквы или скобки"
	}

	return true, "Пароль валиден"
}

func (s *Storage) ValidateGameName(name string) (bool, string) {
	if valid, count := s.isValidLength(name, 3, 100); !valid {
		if count > 0 {
			return false, "Название не может быть больше 100 символов"
		} else if count < 0 {
			return false, "Название не может быть меньше 3 символов"
		}
	}

	if valid := s.isValidString(name, true, true, []rune{'.', ',', '-', '_', '!'}); !valid {
		return false, "Название может содержать только латинские буквы, цифры, точку, запятую, дефис, нижнее подчёркивание или восклицательный знак"
	}

	return true, "Название валидно"
}

func (s *Storage) isValidLength(str string, min, max int) (bool, int) {
	charCount := utf8.RuneCountInString(str)
	if charCount < min {
		return false, charCount - min
	} else if charCount > max {
		return false, charCount - max
	}

	return true, 0
}

func (s *Storage) isValidString(str string, letters bool, numbers bool, additional []rune) bool {
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

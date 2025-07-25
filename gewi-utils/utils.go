package gewiutils

func Ternary(condition bool, pos, neg any) any {
	if condition {
		return pos
	} else {
		return neg
	}
}

func TernaryInt(condition bool, pos, neg int) int {
	if condition {
		return pos
	} else {
		return neg
	}
}

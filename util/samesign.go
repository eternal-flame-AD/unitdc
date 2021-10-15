package util

func SameSign(a int, b int) bool {
	return a > 0 && b > 0 ||
		a < 0 && b < 0
}

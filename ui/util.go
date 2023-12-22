package ui

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// InsertString inserts a string into another string at a given index.
func InsertString(s string, insert string, n int) string {
	return s[:n] + insert + s[n:]
}

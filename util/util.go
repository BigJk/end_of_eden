package util

import (
	"strings"
)

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

// CopyMap copies a map. If the value is a pointer, the pointer is copied, not the value.
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// InsertString inserts a string into another string at a given index.
func InsertString(s string, insert string, n int) string {
	return s[:n] + insert + s[n:]
}

// RemoveAnsiReset removes the first ansi reset code from a string.
func RemoveAnsiReset(s string) string {
	return strings.Replace(s, "\x1b[0m", "", 1)
}

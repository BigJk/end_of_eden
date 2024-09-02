package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// Numbers is a slice of strings that represent the numbers 0-9 in a 5x5 grid.
var Numbers = []string{
	` ██████
██  ████
██ ██ ██
████  ██
 ██████`,
	` ██
███
 ██
 ██
 ██  `,
	`██████
     ██
 █████
██
███████ `,
	`██████
     ██
 █████
     ██
██████  `,
	`██   ██
██   ██
███████
     ██
     ██`,
	`███████
██
███████
     ██
███████`,
	` ██████
██
███████
██    ██
 ██████`,
	`███████
     ██
    ██
   ██
   ██`,
	` █████
██   ██
 █████
██   ██
 █████`,
	` █████
██   ██
 ██████
     ██
 █████`,
}

// GetNumber returns a string representation of a number.
func GetNumber(number int) string {
	return strings.Join(lo.Map([]rune(fmt.Sprint(number)), func(char rune, index int) string {
		num, _ := strconv.Atoi(string(char))
		return Numbers[num]
	}), "")
}

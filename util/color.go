package util

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func RGBColor(r, g, b byte) lipgloss.Color {
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}

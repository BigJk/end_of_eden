package components

import (
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
)

func Error(width int, height int, msg string) string {
	err := lipgloss.NewStyle().Width(width-30).Border(lipgloss.ThickBorder(), true).Padding(0, 2, 0, 1).BorderForeground(style.BaseGray).Foreground(style.BaseWhite).Render(msg)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, err, lipgloss.WithWhitespaceChars("!"), lipgloss.WithWhitespaceForeground(style.BaseGrayDarker))
}

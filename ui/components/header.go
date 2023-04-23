package components

import (
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerLogo = lipgloss.NewStyle().Foreground(style.BaseRedDarker).Render(`▀▄.▀·██▪ ██ ▀▄.▀·•█▌▐█
▐▀▀▪▄▐█· ▐█▌▐▀▀▪▄▐█▐▐▌
▐█▄▄▌██. ██ ▐█▄▄▌██▐█▌`)
	headerOuterStyle = lipgloss.NewStyle().Foreground(style.BaseWhite).Border(lipgloss.BlockBorder(), false, false, true, false).BorderForeground(style.BaseRedDarker)
)

type HeaderValue struct {
	Text  string
	Color lipgloss.Color
}

func NewHeaderValue(text string, color lipgloss.Color) HeaderValue {
	return HeaderValue{text, color}
}

func Header(width int, values []HeaderValue, desc string, other ...string) string {
	entries := []string{
		headerLogo,
	}

	for i := range values {
		entries = append(entries, lipgloss.NewStyle().Bold(true).Foreground(values[i].Color).Padding(0, 3, 0, 3).Render(values[i].Text))
	}

	if len(desc) > 0 {
		entries = append(entries, lipgloss.NewStyle().Italic(true).Foreground(style.BaseGray).Padding(0, 3, 0, 3).Render("\""+desc+"\""))
	}

	entries = append(entries, other...)

	return headerOuterStyle.Width(width).Render(lipgloss.JoinHorizontal(lipgloss.Center, entries...))
}

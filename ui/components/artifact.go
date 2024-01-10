package components

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	artifactStyle = lipgloss.NewStyle().Padding(1, 2).Margin(0, 2)
)

func ArtifactCard(session *game.Session, guid string, baseHeight int, maxHeight int, optionalWidth ...int) string {
	art, _ := session.GetArtifact(guid)
	width := 30
	if len(optionalWidth) > 0 {
		width = optionalWidth[0]
	}

	artifactStyle := artifactStyle.Copy().
		Width(width).
		Border(lipgloss.ThickBorder(), true, false, false, false).
		BorderBackground(lipgloss.Color("#495057")).
		BorderForeground(lipgloss.Color("#495057")).
		Background(lipgloss.Color("#343a40")).
		Foreground(style.BaseWhite)

	tagsText := strings.Join(art.Tags, ", ")

	return artifactStyle.
		Height(baseHeight).
		Render(fmt.Sprintf("%s\n\n%s\n\n%s", style.BoldStyle.Render(art.Name, strings.Repeat(" ", ui.Max(width-6-lipgloss.Width(art.Name)-lipgloss.Width(tagsText), 0)), tagsText), art.Description, lipgloss.NewStyle().Bold(true).Foreground(style.BaseYellow).Render(fmt.Sprintf("%d$", art.Price))))

}

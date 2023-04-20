package components

import (
	"fmt"
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/ui/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	artifactStyle = lipgloss.NewStyle().Padding(1, 2).Margin(0, 2)
)

func ArtifactCard(session *game.Session, guid string, baseHeight int, maxHeight int) string {
	art, _ := session.GetArtifact(guid)

	artifactStyle := artifactStyle.Copy().
		Width(30).
		Border(lipgloss.ThickBorder(), true, false, false, false).
		BorderBackground(lipgloss.Color("#495057")).
		BorderForeground(lipgloss.Color("#495057")).
		Background(lipgloss.Color("#343a40")).
		Foreground(style.BaseWhite)

	return artifactStyle.
		Height(baseHeight).
		Render(fmt.Sprintf("%s\n\n%s\n\n%s", style.BoldStyle.Render(art.Name), art.Description, lipgloss.NewStyle().Bold(true).Foreground(style.BaseYellow).Render(fmt.Sprintf("%d$", art.Price))))

}

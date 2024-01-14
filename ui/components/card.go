package components

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"
	"strings"
)

var (
	cardStyle     = lipgloss.NewStyle().Padding(1, 2).Margin(0, 2)
	headerStlye   = lipgloss.NewStyle().Bold(true)
	cantCastStyle = lipgloss.NewStyle().Foreground(style.BaseRed)
)

func HalfCard(session *game.Session, guid string, active bool, baseHeight int, maxHeight int, minimal bool, width int, checkCasting bool) string {
	fight := session.GetFight()
	card, _ := session.GetCard(guid)
	canCast := !checkCasting || fight.CurrentPoints >= card.PointCost
	cardState := session.GetCardState(guid)

	pointText := strings.Repeat("•", card.PointCost)
	tagsText := strings.Join(card.PublicTags(), ", ")

	cardCol, _ := colorful.Hex(card.Color)
	bgCol, _ := colorful.MakeColor(style.BaseGrayDarker)

	if width <= 0 {
		width = 30
	}

	cardStyle := cardStyle.Copy().
		Width(lo.Ternary(minimal && !active, 10, width)).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderBackground(lipgloss.Color(card.Color)).
		BorderForeground(lo.Ternary(active, style.BaseGray, lipgloss.Color(card.Color))).
		Background(lipgloss.Color(cardCol.BlendRgb(bgCol, 0.6).Hex())).
		Foreground(style.BaseWhite)

	header := headerStlye.Render(fmt.Sprintf("%s%s%s", pointText, strings.Repeat(" ", ui.Max(width-4-lipgloss.Width(pointText)-lipgloss.Width(tagsText), 0)), tagsText))
	if !canCast {
		header = cantCastStyle.Render(header)
	}

	if active {
		return cardStyle.
			Height(ui.Min(maxHeight-1, baseHeight+5)).
			Render(fmt.Sprintf("%s\n\n%s\n\n%s", header, style.BoldStyle.Render(card.Name), cardState))
	}

	if minimal {
		return cardStyle.
			Height(baseHeight).
			Render(fmt.Sprintf("%s\n\n%s", pointText, style.BoldStyle.Render(strings.Join(lo.ChunkString(card.Name, 1), "\n"))))
	}

	return cardStyle.
		Height(baseHeight).
		Render(fmt.Sprintf("%s\n\n%s\n\n%s", header, style.BoldStyle.Render(card.Name), cardState))

}

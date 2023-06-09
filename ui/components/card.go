package components

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/BigJk/end_of_eden/util"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"
	"strings"
)

var (
	cardStyle     = lipgloss.NewStyle().Padding(1, 2).Margin(0, 2)
	cantCastStyle = lipgloss.NewStyle().Foreground(style.BaseRed)
)

func HalfCard(session *game.Session, guid string, active bool, baseHeight int, maxHeight int, minimal bool) string {
	fight := session.GetFight()
	card, _ := session.GetCard(guid)
	canCast := fight.CurrentPoints >= card.PointCost
	cardState := session.GetCardState(guid)

	pointText := strings.Repeat("•", card.PointCost)
	if !canCast {
		pointText = cantCastStyle.Render(pointText)
	}

	cardCol, _ := colorful.Hex(card.Color)
	bgCol, _ := colorful.MakeColor(style.BaseGrayDarker)

	cardStyle := cardStyle.Copy().
		Width(lo.Ternary(minimal && !active, 10, 30)).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderBackground(lipgloss.Color(card.Color)).
		BorderForeground(lo.Ternary(active, style.BaseGray, lipgloss.Color(card.Color))).
		Background(lipgloss.Color(cardCol.BlendRgb(bgCol, 0.6).Hex())).
		Foreground(style.BaseWhite)

	if active {
		return cardStyle.
			Height(util.Min(maxHeight-1, baseHeight+5)).
			Render(fmt.Sprintf("%s\n\n%s\n\n%s", pointText, style.BoldStyle.Render(card.Name), cardState))
	}

	if minimal {
		return cardStyle.
			Height(baseHeight).
			Render(fmt.Sprintf("%s\n\n%s", pointText, style.BoldStyle.Render(strings.Join(lo.ChunkString(card.Name, 1), "\n"))))
	}

	return cardStyle.
		Height(baseHeight).
		Render(fmt.Sprintf("%s\n\n%s\n\n%s", pointText, style.BoldStyle.Render(card.Name), cardState))

}

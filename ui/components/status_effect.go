package components

import (
	"fmt"
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/util"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/samber/lo"
	"strings"
)

func StatusEffect(session *game.Session, guid string) string {
	status := session.GetStatusEffect(guid)
	if status == nil {
		return ""
	}

	fg, _ := colorful.Hex(status.Foreground)
	bg := fg.BlendRgb(colorful.LinearRgb(0, 0, 0), 0.7)

	return fmt.Sprint(session.GetInstance(guid).(game.StatusEffectInstance).Stacks) + lipgloss.NewStyle().Foreground(lipgloss.Color(status.Foreground)).Background(lipgloss.Color(bg.Hex())).Render(status.Look)
}

func StatusEffects(session *game.Session, actor game.Actor) string {
	return strings.Join(lo.Map(util.SortStringsStable(actor.StatusEffects.ToSlice()), func(guid string, index int) string {
		status := session.GetStatusEffect(guid)
		if status == nil {
			return ""
		}

		fg, _ := colorful.Hex(status.Foreground)
		bg := fg.BlendRgb(colorful.LinearRgb(0, 0, 0), 0.7)

		return fmt.Sprint(session.GetInstance(guid).(game.StatusEffectInstance).Stacks) + lipgloss.NewStyle().Foreground(lipgloss.Color(status.Foreground)).Background(lipgloss.Color(bg.Hex())).Render(status.Look)
	}), " ")
}

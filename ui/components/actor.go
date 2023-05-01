package components

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

var actorFaceStyle = lipgloss.NewStyle().Border(lipgloss.OuterHalfBlockBorder()).Padding(0, 1).Margin(0, 0, 1, 0).BorderForeground(style.BaseGrayDarker).Foreground(style.BaseRed)

func Actor(session *game.Session, actor game.Actor, enemy *game.Enemy, showStatus bool, showHp bool, active bool, additional ...string) string {
	face := actorFaceStyle.Copy().BorderForeground(lo.Ternary(active, style.BaseWhite, style.BaseGrayDarker)).Foreground(lipgloss.Color(enemy.Color)).Render(enemy.Look)

	var parts []string

	if showStatus {
		parts = append(parts, StatusEffects(session, actor)+"\n")
	}

	parts = append(parts, face, enemy.Name)

	if showHp {
		parts = append(parts, fmt.Sprintf("%d / %d", actor.HP, enemy.MaxHP))
	}

	parts = append(parts, lo.Filter(additional, func(item string, index int) bool {
		return len(item) > 0
	})...)

	return lipgloss.NewStyle().Foreground(style.BaseWhite).Margin(0, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Center, parts...,
		))
}

package gameview

import (
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/menus"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type DeathAnimationFrame time.Time

type DeathAnimationModel struct {
	width    int
	height   int
	target   game.Actor
	death    game.StateEventDeathData
	progress float64
}

func NewDeathAnimationModel(width int, height int, target game.Actor, death game.StateEventDeathData) DeathAnimationModel {
	return DeathAnimationModel{
		width:  width,
		height: height,
		target: target,
		death:  death,
	}
}

func (m DeathAnimationModel) SetSize(width int, height int) DeathAnimationModel {
	m.width = width
	m.height = height
	return m
}

func (m DeathAnimationModel) Init() tea.Cmd {
	return nil
}

func (m DeathAnimationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case DeathAnimationFrame:
		elapsed := (1.0 / 30.0) / 5.0
		m.progress += elapsed
	}

	if m.progress >= 1.0 {
		return nil, nil
	}

	return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
		return DeathAnimationFrame(t)
	})
}

func (m DeathAnimationModel) View() string {
	headerStyle := lipgloss.NewStyle().Margin(4).Padding(2).Border(lipgloss.NormalBorder())

	faceStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(menus.BaseWhite).Foreground(menus.BaseRed)
	face := faceStyle.Render("@")

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			headerStyle.Render("Enemy Slayed"),
			face,
			m.target.Name,
		),
	)
}

package gameview

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/ui/animation"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"time"
)

type DeathAnimationFrame string

type DeathAnimationModel struct {
	id       string
	width    int
	height   int
	target   game.Actor
	enemy    *game.Enemy
	death    game.StateEventDeathData
	progress float64
	started  bool
}

func NewDeathAnimationModel(width int, height int, target game.Actor, targetEnemy *game.Enemy, death game.StateEventDeathData) DeathAnimationModel {
	return DeathAnimationModel{
		id:     fmt.Sprint(rand.Intn(100000)),
		width:  width,
		height: height,
		target: target,
		enemy:  targetEnemy,
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height)
	case tea.Key:
		if m.progress > 0.1 && (msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace) {
			return nil, nil
		}
	case tea.MouseMsg:
		if m.progress > 0.1 && (msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft) {
			return nil, nil
		}
	case DeathAnimationFrame:
		if string(msg) != m.id {
			return m, nil
		}

		if m.progress == 0 {
			audio.Play("death_scream_1")
		}

		elapsed := 1.0 / 30.0 / 5.0
		m.progress += elapsed

		if m.progress >= 1.0 {
			return nil, nil
		}

		return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
			return DeathAnimationFrame(m.id)
		})
	}

	if !m.started {
		m.started = true
		return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
			return DeathAnimationFrame(m.id)
		})
	}

	return m, nil
}

const killedText = `▄ •▄ ▪  ▄▄▌  ▄▄▌  ▄▄▄ .·▄▄▄▄  ▄▄ ▄▄ 
█▌▄▌▪██ ██•  ██•  ▀▄.▀·██▪ ██ ██▌██▌
▐▀▀▄·▐█·██▪  ██▪  ▐▀▀▪▄▐█· ▐█▌▐█·▐█·
▐█.█▌▐█▌▐█▌▐▌▐█▌▐▌▐█▄▄▌██. ██ .▀ .▀ 
·▀  ▀▀▀▀.▀▀▀ .▀▀▀  ▀▀▀ ▀▀▀▀▀•  ▀  ▀`

func (m DeathAnimationModel) View() string {
	killedText := animation.JitterText(killedText, m.progress*2.5, 0, 10)

	face := faceStyle.Copy().BorderForeground(style.BaseRed).Foreground(lipgloss.Color(m.enemy.Color)).Render(m.enemy.Look)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			style.RedText.Render(`━╋━
 ┃`),
			face,
			style.BaseText.Render(m.target.Name),
			lipgloss.NewStyle().Margin(2, 0, 0, 0).Foreground(style.BaseRed).Render(killedText),
		),
	)
}

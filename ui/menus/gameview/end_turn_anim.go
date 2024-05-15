package gameview

import (
	"fmt"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"time"
)

var round = `▄▄▄        ▄• ▄▌ ▐ ▄ ·▄▄▄▄  
▀▄ █·▪     █▪██▌•█▌▐███▪ ██ 
▐▀▀▄  ▄█▀▄ █▌▐█▌▐█▐▐▌▐█· ▐█▌
▐█•█▌▐█▌.▐▌▐█▄█▌██▐█▌██. ██ 
.▀  ▀ ▀█▄▀▪ ▀▀▀ ▀▀ █▪▀▀▀▀▀• `

type EndTurnAnimationFrame string

type EndTurnAnimationModel struct {
	id     string
	width  int
	height int

	started bool
	elapsed float64
	turn    int
}

func NewEndTurnAnimationModel(width int, height int, turn int) EndTurnAnimationModel {
	return EndTurnAnimationModel{
		id:     fmt.Sprint(rand.Intn(100000)),
		width:  width,
		height: height,
		turn:   turn,
	}
}

func (m EndTurnAnimationModel) SetSize(width int, height int) EndTurnAnimationModel {
	m.width = width
	m.height = height
	return m
}

func (m EndTurnAnimationModel) Init() tea.Cmd {
	return nil
}

func (m EndTurnAnimationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height)
	case tea.Key:
		if m.elapsed > 0.1 && (msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace) {
			return nil, nil
		}
	case tea.MouseMsg:
		if m.elapsed > 0.1 && (msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft) {
			return nil, nil
		}
	case EndTurnAnimationFrame:
		if string(msg) == m.id {
			m.elapsed += 1.0 / 30.0

			if m.elapsed > 0.5 {
				return nil, nil
			}

			return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
				return EndTurnAnimationFrame(m.id)
			})
		}
	}

	// Send first tick
	if !m.started {
		m.started = true
		return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
			return EndTurnAnimationFrame(m.id)
		})
	}

	return m, nil
}

func (m EndTurnAnimationModel) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(style.BaseRed).Margin(0, 0, 3, 0).Render(round),
		style.RedText.Render(ui.GetNumber(m.turn)),
	))
}

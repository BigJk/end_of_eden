package gameover

import (
	"fmt"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/animation"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"time"
)

const text = ` ▄▄ •  ▄▄▄· • ▌ ▄ ·. ▄▄▄ .           ▌ ▐·▄▄▄ .▄▄▄           
▐█ ▀ ▪▐█ ▀█ ·██ ▐███▪▀▄.▀·    ▪     ▪█·█▌▀▄.▀·▀▄ █·         
▄█ ▀█▄▄█▀▀█ ▐█ ▌▐▌▐█·▐▀▀▪▄     ▄█▀▄ ▐█▐█•▐▀▀▪▄▐▀▀▄          
▐█▄▪▐█▐█ ▪▐▌██ ██▌▐█▌▐█▄▄▌    ▐█▌.▐▌ ███ ▐█▄▄▌▐█•█▌         
·▀▀▀▀  ▀  ▀ ▀▀  █▪▀▀▀ ▀▀▀      ▀█▄▀▪. ▀   ▀▀▀ .▀  ▀ ▀  ▀  ▀ `

const (
	ZoneToMenu = "to_menu"
)

type GameOverFrame time.Time

type Model struct {
	ui.MenuBase

	zones     *zone.Manager
	started   bool
	progress  float64
	lastMouse tea.MouseMsg

	Session *game.Session
	Start   game.StateCheckpointMarker

	allDamage         int
	allDamageReceived int
	allGold           int
}

func New(zones *zone.Manager, session *game.Session, start game.StateCheckpointMarker) Model {
	m := Model{
		zones:   zones,
		Session: session,
		Start:   start,
	}

	// Collect stats
	diff := start.Diff(session)
	for i := range diff {
		if val, ok := diff[i].Events[game.StateEventDamage]; ok {
			dmg := val.(game.StateEventDamageData)
			if dmg.Target != game.PlayerActorID {
				m.allDamage += dmg.Damage
			} else {
				m.allDamageReceived += dmg.Damage
			}
		}

		if val, ok := diff[i].Events[game.StateEventMoney]; ok {
			mon := val.(game.StateEventMoneyData)
			m.allGold += mon.Money
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case tea.MouseMsg:
		m.lastMouse = msg

		if msg.Type == tea.MouseLeft && m.zones.Get(ZoneToMenu).InBounds(msg) {
			m.Session.Close()
			return nil, nil
		}
	case GameOverFrame:
		if m.progress == 0 {
			audio.Play("game_over")
		}

		elapsed := 1.0 / 30.0 / 1.5
		m.progress += elapsed

		if m.progress >= 1.0 {
			m.progress = 1.0
		} else {
			return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
				return GameOverFrame(t)
			})
		}
	}

	if !m.started {
		m.started = true
		return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
			return GameOverFrame(t)
		})
	}

	return m, nil
}

func (m Model) View() string {
	top := m.top()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		top,
		lipgloss.Place(m.Size.Width, m.Size.Height-lipgloss.Height(top), lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center,
				style.RedText.Render(animation.JitterText(text, m.progress, 0, 10)),
				lipgloss.NewStyle().Margin(2, 0, 1, 0).Padding(1, 3).Border(lipgloss.ThickBorder()).BorderForeground(style.BaseRedDarker).Foreground(style.BaseWhite).Render(
					fmt.Sprintf(
						"%s\n\n%s%d\n%s%d\n%s%d\n%s%d",
						style.BoldStyle.Render("Run Statistic"),
						style.BoldStyle.Render(fmt.Sprintf("%-20s :  ", "Stages ")), m.Session.GetStagesCleared(),
						style.BoldStyle.Render(fmt.Sprintf("%-20s :  ", "Damage Done ")), m.allDamage,
						style.BoldStyle.Render(fmt.Sprintf("%-20s :  ", "Damage Received ")), m.allDamageReceived,
						style.BoldStyle.Render(fmt.Sprintf("%-20s :  ", "Gold Collected ")), m.allGold,
					),
				),
				m.zones.Mark(ZoneToMenu, style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneToMenu).InBounds(m.lastMouse), style.BaseRed, style.BaseRedDarker)).Render("Accept your fate...")),
			),
		),
	)
}

func (m Model) top() string {
	outerStyle := lipgloss.NewStyle().
		Width(m.Size.Width).
		Foreground(style.BaseWhite).
		Border(lipgloss.BlockBorder(), false, false, true, false).
		BorderForeground(style.BaseRedDarker)

	fight := m.Session.GetFight()
	player := m.Session.GetPlayer()

	return outerStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(style.BaseRedDarker).Render(`▐█ ▀ ▪▪     •█▌▐█▪▀·.█▌▪     
▄█ ▀█▄ ▄█▀▄ ▐█▐▐▌▄█▀▀▀• ▄█▀▄ 
▐█▄▪▐█▐█▌.▐▌██▐█▌█▌▪▄█▀▐█▌.▐▌`),
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00")).Padding(0, 4, 0, 4).Render(fmt.Sprintf("Gold: %d", player.Gold)),
		lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Padding(0, 4, 0, 0).Render(fmt.Sprintf("HP: %d / %d", player.HP, player.MaxHP)),
		lipgloss.NewStyle().Bold(true).Foreground(style.BaseWhite).Padding(0, 4, 0, 0).Render(fmt.Sprintf("%d. Stage", m.Session.GetStagesCleared()+1)),
		lipgloss.NewStyle().Bold(true).Foreground(style.BaseWhite).Padding(0, 4, 0, 0).Render(fmt.Sprintf("%d. Round", fight.Round+1)),
		lipgloss.NewStyle().Italic(true).Foreground(style.BaseGray).Padding(0, 4, 0, 0).Render("\"Better luck next time...\""),
	))
}

package carousel

import (
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"strings"
)

const (
	ZoneLeftButton  = "left_button"
	ZoneRightButton = "right_button"
	ZoneDoneButton  = "done_button"
)

type Model struct {
	ui.MenuBase

	parent   tea.Model
	title    string
	items    []string
	selected int

	onceFn func()
}

func New(parent tea.Model, zones *zone.Manager, title string, items []string) Model {
	return Model{
		MenuBase: ui.NewMenuBase().WithZones(zones),
		parent:   parent,
		title:    title,
		items:    items,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			audio.Play("btn_menu")
			return m.parent, nil
		} else if msg.Type == tea.KeyLeft {
			if m.selected > 0 {
				m.selected--
				audio.Play("btn_menu")
			}
		} else if msg.Type == tea.KeyRight {
			if m.selected < len(m.items)-1 {
				m.selected++
				audio.Play("btn_menu")
			}
		}
	case tea.MouseMsg:
		m.LastMouse = msg
		if msg.Button != tea.MouseButtonNone {
			if m.ZoneInBounds(ZoneLeftButton) {
				if m.selected > 0 {
					m.selected--
					audio.Play("btn_menu")
				}
			}
			if m.ZoneInBounds(ZoneRightButton) {
				if m.selected < len(m.items)-1 {
					m.selected++
					audio.Play("btn_menu")
				}
			}
			if m.ZoneInBounds(ZoneDoneButton) {
				audio.Play("btn_menu")
				return m.parent, nil
			}
		}
	}

	return m, nil
}

func (m Model) leftButton() string {
	background := m.ZoneBackground(ZoneLeftButton, style.BaseRed, style.BaseRedDarker)
	if m.selected == 0 {
		background = style.BaseGrayDarker
	}

	return m.ZoneMark(ZoneLeftButton, style.HeaderStyle.Copy().Background(background).Margin(0, 2).Render("<--"))
}

func (m Model) rightButton() string {
	background := m.ZoneBackground(ZoneRightButton, style.BaseRed, style.BaseRedDarker)
	if m.selected == len(m.items)-1 {
		background = style.BaseGrayDarker
	}

	return m.ZoneMark(ZoneRightButton, style.HeaderStyle.Copy().Background(background).Margin(0, 2).Render("-->"))
}

func (m Model) View() string {
	title := style.BoldStyle.Copy().MarginBottom(4).Render(m.title)

	middle := lipgloss.JoinHorizontal(lipgloss.Center,
		m.leftButton(),
		m.items[m.selected],
		m.rightButton(),
	)

	dots := lipgloss.NewStyle().Margin(2, 0).Render(strings.Join(lo.Map(m.items, func(item string, index int) string {
		if index == m.selected {
			return "●"
		}
		return "○"
	}), " "))

	doneButton := style.HeaderStyle.Copy().Background(m.ZoneBackground(ZoneDoneButton, style.BaseRed, style.BaseRedDarker)).Render(m.ZoneMark(ZoneDoneButton, "Continue"))

	return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, title, middle, dots, doneButton))
}

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

	zones     *zone.Manager
	parent    tea.Model
	lastMouse tea.MouseMsg
	title     string
	items     []string
	selected  int

	onceFn func()
}

func New(parent tea.Model, zones *zone.Manager, title string, items []string) Model {
	return Model{
		zones:  zones,
		parent: parent,
		title:  title,
		items:  items,
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
		if msg.Type == tea.KeyEscape {
			return m.parent, nil
		} else if msg.Type == tea.KeyEnter {
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
		if msg.Type == tea.MouseLeft {
			if m.zones.Get(ZoneLeftButton).InBounds(msg) {
				if m.selected > 0 {
					m.selected--
					audio.Play("btn_menu")
				}
			}
			if m.zones.Get(ZoneLeftButton).InBounds(msg) {
				if m.selected < len(m.items)-1 {
					m.selected++
					audio.Play("btn_menu")
				}
			}
			if m.zones.Get(ZoneDoneButton).InBounds(msg) {
				audio.Play("btn_menu")
				return m.parent, nil
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	title := style.BoldStyle.Copy().MarginBottom(4).Render(m.title)

	leftButton := m.zones.Mark(ZoneLeftButton, style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneLeftButton).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2).Render("<--"))
	rightButton := m.zones.Mark(ZoneRightButton, style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneRightButton).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2).Render("-->"))
	middle := lipgloss.JoinHorizontal(lipgloss.Center,
		leftButton,
		m.items[m.selected],
		rightButton,
	)

	dots := lipgloss.NewStyle().Margin(2, 0).Render(strings.Join(lo.Map(m.items, func(item string, index int) string {
		if index == m.selected {
			return "●"
		}
		return "○"
	}), " "))

	doneButton := m.zones.Mark(ZoneDoneButton, style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneDoneButton).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).MarginTop(2).Render("Continue"))

	return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, title, middle, dots, doneButton))
}

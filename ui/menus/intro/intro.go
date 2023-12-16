package intro

import (
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/components/gifviewer"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var frameStyle = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder(), true).
	BorderForeground(style.BaseRedDarker).
	Padding(0, 1)

type Model struct {
	ui.MenuBase

	parent tea.Model
	gif    tea.Model
}

func New(parent tea.Model) Model {
	return Model{
		parent: parent,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gif == nil {
		gif, err := gifviewer.New(m, "intro.gif", 10, 100, 0)
		if err != nil {
			return m.parent, nil
		}
		m.gif = gif
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case tea.KeyMsg:
		if msg.Type == tea.KeyEscape {
			return m.parent, nil
		}
	}

	var cmd tea.Cmd
	m.gif, cmd = m.gif.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.gif == nil {
		return ""
	}

	frame := m.gif.View()
	frameFramed := frameStyle.Render(
		lipgloss.NewStyle().MaxWidth(100 - 4).Render(frame[:len(frame)-1]))

	return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(
		lipgloss.Center,
		frameFramed,
		frameStyle.Render(lipgloss.NewStyle().Width(100-4).Render(ui.About)),
	))
}

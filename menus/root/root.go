package root

import (
	"github.com/BigJk/project_gonzo/menus"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	current tea.Model
	size    tea.WindowSizeMsg
}

func New(root tea.Model) Model {
	return Model{
		current: root,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)

	if menu, ok := m.current.(menus.Menu); ok && !menu.HasSize() {
		return m, tea.Batch(cmd, func() tea.Msg {
			return m.size
		})
	}

	return m, cmd
}

func (m Model) View() string {
	return zone.Scan(m.current.View())
}

package root

import (
	"github.com/BigJk/project_gonzo/ui"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	zones   *zone.Manager
	root    tea.Model
	current tea.Model
	size    tea.WindowSizeMsg
}

func New(zones *zone.Manager, root tea.Model) Model {
	return Model{
		zones:   zones,
		root:    root,
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

	if menu, ok := m.current.(ui.Menu); ok && !menu.HasSize() {
		return m, tea.Batch(cmd, func() tea.Msg {
			return m.size
		})
	}

	if m.current == nil {
		// Fall back to main menu
		m.current = m.root
	}

	return m, cmd
}

func (m Model) View() string {
	return m.zones.Scan(m.current.View())
}

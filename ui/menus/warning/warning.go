package warning

import (
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ui.MenuBase

	parent tea.Model
	text   string
}

func New(parent tea.Model, text string) Model {
	return Model{parent: parent, text: text}
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
		}
	}

	return m, nil
}

func (m Model) View() string {
	return components.Error(m.Size.Width, m.Size.Height, m.text)
}

package root

import tea "github.com/charmbracelet/bubbletea"

type OnVisibleModel struct {
	parent tea.Model
	fn     func(model tea.Model)
}

func NewOnVisibleModel(parent tea.Model, fn func(model tea.Model)) OnVisibleModel {
	return OnVisibleModel{parent: parent, fn: fn}
}

func (m OnVisibleModel) Init() tea.Cmd {
	return nil
}

func (m OnVisibleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.parent == nil {
		return nil, nil
	}

	switch msg.(type) {
	case ModelGettingVisibleMsg:
		m.fn(m.parent)
		return m.parent.Update(msg)
	}

	pm, cmd := m.parent.Update(msg)
	m.parent = pm

	return m, cmd
}

func (m OnVisibleModel) View() string {
	if m.parent == nil {
		return ""
	}
	return m.parent.View()
}

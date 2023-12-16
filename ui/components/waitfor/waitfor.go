package waitfor

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	model tea.Model
	cond  func(msg tea.Msg) bool
}

func New(model tea.Model, cond func(msg tea.Msg) bool) tea.Model {
	return Model{
		model: model,
		cond:  cond,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.cond(msg) {
		return nil, nil
	}

	model, modelCmd := m.model.Update(msg)
	if model == nil {
		return nil, nil
	}
	m.model = model

	return m, modelCmd
}

func (m Model) View() string {
	return m.model.View()
}

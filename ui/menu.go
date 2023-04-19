package ui

import tea "github.com/charmbracelet/bubbletea"

// Menu is a tea.Model that keeps track of its size. It is intended to
// trigger the re-distribution of the tea.WindowSizeMsg for nested models.
type Menu interface {
	tea.Model
	HasSize() bool
}

// MenuBase is the base Menu implementation.
type MenuBase struct {
	Size tea.WindowSizeMsg
}

func (m MenuBase) HasSize() bool {
	return m.Size.Width > 0
}

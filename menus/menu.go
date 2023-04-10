package menus

import tea "github.com/charmbracelet/bubbletea"

type Menu interface {
	tea.Model
	HasSize() bool
}

type MenuBase struct {
	Size tea.WindowSizeMsg
}

func (m MenuBase) HasSize() bool {
	return m.Size.Width > 0
}

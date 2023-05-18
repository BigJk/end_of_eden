package root

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PushModelMsg tea.Model

// Push pushes a new model on the root stack.
func Push(model tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushModelMsg(model)
	}
}

// Tooltip represents a tooltip aka overlay message  that should be displayed.
type Tooltip struct {
	ID      string
	Content string
	X       int
	Y       int
}

type TooltipMsg Tooltip

// TooltipCreate creates a new tooltip.
func TooltipCreate(tip Tooltip) tea.Cmd {
	return func() tea.Msg {
		return TooltipMsg(tip)
	}
}

type TooltipDeleteMsg string

// TooltipDelete deletes a tooltip.
func TooltipDelete(id string) tea.Cmd {
	return func() tea.Msg {
		return TooltipDeleteMsg(id)
	}
}

type TooltipClearMsg struct{}

// TooltipClear clears all tooltips.
func TooltipClear() tea.Cmd {
	return func() tea.Msg {
		return TooltipClearMsg(struct{}{})
	}
}

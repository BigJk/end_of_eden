package root

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PushModelMsg []tea.Model

// Push pushes a new model on the root stack.
func Push(model tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushModelMsg([]tea.Model{model})
	}
}

// PushAll pushes multiple models on the root stack.
func PushAll(models ...tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushModelMsg(models)
	}
}

type ModelGettingVisibleMsg struct{}

// GettingVisible is a message that is sent to a model when it is getting visible.
func GettingVisible() tea.Cmd {
	return func() tea.Msg {
		return ModelGettingVisibleMsg(struct{}{})
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

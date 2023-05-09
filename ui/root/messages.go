package root

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PushModelMsg tea.Model

func Push(model tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushModelMsg(model)
	}
}

type Tooltip struct {
	ID      string
	Content string
	X       int
	Y       int
}

type TooltipMsg Tooltip

func TooltipCreate(tip Tooltip) tea.Cmd {
	return func() tea.Msg {
		return TooltipMsg(tip)
	}
}

type TooltipDeleteMsg string

func TooltipDelete(id string) tea.Cmd {
	return func() tea.Msg {
		return TooltipDeleteMsg(id)
	}
}

type TooltipClearMsg struct{}

func TooltipClear() tea.Cmd {
	return func() tea.Msg {
		return TooltipClearMsg(struct{}{})
	}
}

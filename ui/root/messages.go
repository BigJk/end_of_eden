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

type PushTransitionFuncMsg func(parent tea.Model) tea.Model

// PushTransitionFunc pushes a new transition model on the root ui that will be shown between models on the stack.
func PushTransitionFunc(fn func(parent tea.Model) tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushTransitionFuncMsg(fn)
	}
}

// RemovePushTransitionFunc removes the transition model from the root ui.
func RemovePushTransitionFunc() tea.Cmd {
	return func() tea.Msg {
		return PushTransitionFuncMsg(nil)
	}
}

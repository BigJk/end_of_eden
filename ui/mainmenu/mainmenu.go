package mainmenu

import (
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/ui"
	"github.com/BigJk/project_gonzo/ui/about"
	"github.com/BigJk/project_gonzo/ui/gameview"
	"github.com/BigJk/project_gonzo/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Model struct {
	ui.MenuBase

	choices ChoicesModel
}

func NewModel() Model {
	model := Model{
		choices: NewChoicesModel(),
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
		updated, cmd := m.choices.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - (strings.Count(ui.Title, "\n") + style.TitleStyle.GetVerticalFrameSize() + 1),
		})
		m.choices = updated.(ChoicesModel)
		return m, cmd
	}

	updated, cmd := m.choices.Update(msg)
	m.choices = updated.(ChoicesModel)

	switch m.choices.selected {
	case ChoiceContinue:
	case ChoiceNewGame:
		m.choices = m.choices.Clear()
		return gameview.New(m, game.NewSession(game.WithDebugEnabled("127.0.0.1:8272"))), cmd
	case ChoiceAbout:
		m.choices = m.choices.Clear()
		return about.New(m), cmd
	case ChoiceSettings:
	}

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, style.TitleStyle.Render(ui.Title), m.choices.View())
}
